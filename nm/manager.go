package nm

import (
	"errors"
	"launchpad.net/go-dbus/v1"
	"sync"
)

type Manager struct {
	nm           *dbus.ObjectProxy
	propsGuard   *sync.RWMutex
	props        map[string]dbus.Variant
	propsChanged *dbus.SignalWatch
}

func NewManager(conn *dbus.Connection) (*Manager, error) {
	if nm := conn.Object(Service, Object); nm != nil {
		props := &dbus.Properties{nm}
		all, err := props.GetAll(Interface)

		if err != nil {
			return nil, errors.New("Failed to query properties")
		}

		pcWatch, err := nm.WatchSignal(Interface, "org.freedesktop.Properties.PropertiesChanged")

		if err != nil {
			return nil, errors.New("Failed to setup a SignalWatch for PropertiesChanged")
		}

		manager := &Manager{nm, &sync.RWMutex{}, all, pcWatch}

		go func() {
			msg := <-manager.propsChanged.C
			changed := make(map[string]dbus.Variant)

			msg.Args(nil, changed, nil)

			manager.propsGuard.Lock()
			defer manager.propsGuard.Unlock()

			for k, v := range changed {
				manager.props[k] = v
			}
		}()

		return manager, nil
	}

	return nil, errors.New("Failed to acquire an ObjectProxy for org.freedesktop.NetworkManager")
}
