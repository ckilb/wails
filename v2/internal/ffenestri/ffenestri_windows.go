package ffenestri

import "C"

/*

#cgo windows CXXFLAGS: -std=c++11
#cgo windows,amd64 LDFLAGS: -lgdi32 -lole32 -lShlwapi -luser32 -loleaut32 -ldwmapi

#include "ffenestri.h"

extern void DisableWindowIcon(struct Application* app);

*/
import "C"

func (a *Application) processPlatformSettings() error {

	config := a.config.Windows
	if config == nil {
		return nil
	}

	// Check if the webview should be transparent
	if config.WebviewIsTransparent {
		C.WebviewIsTransparent(a.app)
	}

	if config.WindowBackgroundIsTranslucent {
		C.WindowBackgroundIsTranslucent(a.app)
	}

	if config.DisableWindowIcon {
		C.DisableWindowIcon(a.app)
	}

	// Process menu
	//applicationMenu := options.GetApplicationMenu(a.config)
	applicationMenu := a.menuManager.GetApplicationMenuJSON()
	println("Appmenu =", applicationMenu)
	if applicationMenu != "" {
		C.SetApplicationMenu(a.app, a.string2CString(applicationMenu))
	}
	//
	//// Process tray
	//trays, err := a.menuManager.GetTrayMenus()
	//if err != nil {
	//	return err
	//}
	//if trays != nil {
	//	for _, tray := range trays {
	//		C.AddTrayMenu(a.app, a.string2CString(tray))
	//	}
	//}
	//
	//// Process context menus
	//contextMenus, err := a.menuManager.GetContextMenus()
	//if err != nil {
	//	return err
	//}
	//if contextMenus != nil {
	//	for _, contextMenu := range contextMenus {
	//		C.AddContextMenu(a.app, a.string2CString(contextMenu))
	//	}
	//}
	//
	//// Process URL Handlers
	//if a.config.Mac.URLHandlers != nil {
	//	C.HasURLHandlers(a.app)
	//}

	return nil
}