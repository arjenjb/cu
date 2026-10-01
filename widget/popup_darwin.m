#import <AppKit/AppKit.h>
#import <objc/runtime.h>


// Marks the windows that are popups.
static char popupKey;

static BOOL isPopup(NSView *view) {
	return view.window != nil && objc_getAssociatedObject(view.window, &popupKey) != nil;
}

// Deliver the first click in the popup, which is never the key window.
static BOOL popupAcceptsFirstMouse(NSView *self, SEL _cmd, NSEvent *event) {
	return isPopup(self);
}

// Gio has no handler for the mouse leaving the view. Send a move outside of
// the view instead so hovered items are reset.
static void popupMouseExited(NSView *self, SEL _cmd, NSEvent *event) {
	if (!isPopup(self)) {
		[self.nextResponder mouseExited:event];
		return;
	}
	NSEvent *move = [NSEvent mouseEventWithType:NSEventTypeMouseMoved
	                                   location:event.locationInWindow
	                              modifierFlags:event.modifierFlags
	                                  timestamp:event.timestamp
	                               windowNumber:event.windowNumber
	                                    context:nil
	                                eventNumber:0
	                                 clickCount:0
	                                   pressure:0];
	[self mouseMoved:move];
}

// Popups never become the key window, so the window of the editor keeps the
// keyboard focus like with a native menu. Gio shows the window before it can
// be made borderless, which would prevent that by itself.
static IMP origCanBecomeKey;
static IMP origCanBecomeMain;

static BOOL popupCanBecomeKey(NSWindow *self, SEL _cmd) {
	if (objc_getAssociatedObject(self, &popupKey) != nil) {
		return NO;
	}
	return ((BOOL(*)(id, SEL))origCanBecomeKey)(self, _cmd);
}

static BOOL popupCanBecomeMain(NSWindow *self, SEL _cmd) {
	if (objc_getAssociatedObject(self, &popupKey) != nil) {
		return NO;
	}
	return ((BOOL(*)(id, SEL))origCanBecomeMain)(self, _cmd);
}

// Add the popup behavior to the classes of Gio. Changing the class of the
// objects instead breaks key value observing, which AppKit relies on.
static void addPopupMethods(Class viewClass) {
	static BOOL added;
	if (added) {
		return;
	}
	added = YES;

	Method m = class_getInstanceMethod([NSWindow class], @selector(canBecomeKeyWindow));
	origCanBecomeKey = method_setImplementation(m, (IMP)popupCanBecomeKey);
	m = class_getInstanceMethod([NSWindow class], @selector(canBecomeMainWindow));
	origCanBecomeMain = method_setImplementation(m, (IMP)popupCanBecomeMain);

	SEL sel = @selector(acceptsFirstMouse:);
	class_addMethod(viewClass, sel, (IMP)popupAcceptsFirstMouse,
	                method_getTypeEncoding(class_getInstanceMethod([NSView class], sel)));
	sel = @selector(mouseExited:);
	class_addMethod(viewClass, sel, (IMP)popupMouseExited,
	                method_getTypeEncoding(class_getInstanceMethod([NSView class], sel)));
}

// The event monitor and notification observer of each open popup.
static NSMutableDictionary<NSNumber *, NSArray *> *popupObservers;

void cu_popupAttach(uintptr_t viewRef) {
	NSView *view = (__bridge NSView *)(void *)viewRef;
	NSWindow *window = view.window;
	if (window == nil || objc_getAssociatedObject(window, &popupKey) != nil) {
		return;
	}

	objc_setAssociatedObject(window, &popupKey, @YES, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
	addPopupMethods([view class]);

	// Keep the window invisible until it is placed
	window.alphaValue = 0;

	// Non-key windows don't get mouse moves by default, which hovering needs.
	NSTrackingArea *area = [[NSTrackingArea alloc]
		initWithRect:NSZeroRect
		     options:NSTrackingMouseMoved | NSTrackingMouseEnteredAndExited |
		             NSTrackingActiveAlways | NSTrackingInVisibleRect
		       owner:view
		    userInfo:nil];
	[view addTrackingArea:area];
}

void cu_mouseLocation(double *x, double *y) {
	// Core Graphics can be used from any thread. Its coordinates have their
	// origin at the top left of the main screen.
	CGEventRef event = CGEventCreate(NULL);
	CGPoint p = CGEventGetLocation(event);
	CFRelease(event);
	*x = p.x;
	*y = p.y;
}

void cu_popupShow(uintptr_t viewRef, uintptr_t popupID, double x, double y, double flipRight, double flipBottom, double width, double height, double radius) {
	NSView *view = (__bridge NSView *)(void *)viewRef;
	NSWindow *window = view.window;
	if (window == nil) {
		return;
	}

	// Changing the style while the view is being attached deallocates it, so
	// it is done here.
	window.styleMask = NSWindowStyleMaskBorderless;
	window.level = NSPopUpMenuWindowLevel;
	window.opaque = NO;
	window.backgroundColor = NSColor.clearColor;
	window.hasShadow = YES;
	view.layer.cornerRadius = radius;
	view.layer.masksToBounds = YES;

	// The positions are in Core Graphics coordinates, with the origin at the
	// top left of the main screen. AppKit has it at the bottom left.
	CGFloat mainTop = NSMaxY(NSScreen.screens.firstObject.frame);
	NSPoint p = NSMakePoint(x, mainTop - y);

	NSRect visible = NSScreen.screens.firstObject.visibleFrame;
	for (NSScreen *screen in NSScreen.screens) {
		if (NSPointInRect(p, screen.frame)) {
			visible = screen.visibleFrame;
			break;
		}
	}
	CGFloat visibleTop = mainTop - NSMaxY(visible);
	CGFloat visibleBottom = mainTop - NSMinY(visible);

	// Move the popup so its right or bottom edge is at the flip position when
	// it doesn't fit on the screen.
	CGFloat left = x;
	CGFloat top = y;
	if (left + width > NSMaxX(visible)) {
		left = MAX(NSMinX(visible), flipRight - width);
	}
	if (top + height > visibleBottom) {
		top = MAX(visibleTop, flipBottom - height);
	}
	CGFloat bottom = mainTop - (top + height);

	[window setFrame:NSMakeRect(left, bottom, width, height) display:YES];
	window.alphaValue = 1;
	[window invalidateShadow];

	// Dismiss the popup on a click outside of it, on escape or when the
	// application is deactivated. The click and escape are swallowed, like a
	// native menu does.
	__weak NSWindow *weakWindow = window;
	NSEventMask mask = NSEventMaskLeftMouseDown | NSEventMaskRightMouseDown |
	                   NSEventMaskOtherMouseDown | NSEventMaskKeyDown;
	id monitor = [NSEvent addLocalMonitorForEventsMatchingMask:mask handler:^NSEvent *(NSEvent *event) {
		if (event.type == NSEventTypeKeyDown) {
			if (event.keyCode == 53) { // escape
				[weakWindow close];
				return nil;
			}
			return event;
		}
		if (event.window != weakWindow) {
			[weakWindow close];
			return nil;
		}
		return event;
	}];
	id observer = [NSNotificationCenter.defaultCenter
		addObserverForName:NSApplicationDidResignActiveNotification
		            object:nil
		             queue:nil
		        usingBlock:^(NSNotification *note) {
			        [weakWindow close];
		        }];

	if (popupObservers == nil) {
		popupObservers = [NSMutableDictionary dictionary];
	}
	popupObservers[@(popupID)] = @[ monitor, observer ];
}

// Gio closes windows with performClose, which does nothing for a window
// without a close button.
void cu_popupClose(uintptr_t viewRef) {
	NSView *view = (__bridge NSView *)(void *)viewRef;
	[view.window close];
}

void cu_popupCleanup(uintptr_t popupID) {
	dispatch_async(dispatch_get_main_queue(), ^{
		NSArray *observers = popupObservers[@(popupID)];
		if (observers == nil) {
			return;
		}
		[NSEvent removeMonitor:observers[0]];
		[NSNotificationCenter.defaultCenter removeObserver:observers[1]];
		[popupObservers removeObjectForKey:@(popupID)];
	});
}
