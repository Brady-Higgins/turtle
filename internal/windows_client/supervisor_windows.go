package windows_client

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type windows_c struct {
	manager *mgr.Mgr
	service *mgr.Service
}

func New() (*windows_c, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, err
	}
	s, err := m.OpenService("cloudflared")
	if err != nil {
		return nil, err
	}
	return &windows_c{manager: m, service: s}, nil
}

func (w *windows_c) end() error {
	err := w.manager.Disconnect()
	if err != nil {
		return err
	}
	err = w.service.Close()
	return err
}

func (w *windows_c) UpdateServiceStartType() error {
	cfg, err := w.service.Config()
	cfg.StartType = mgr.StartManual
	err = w.service.UpdateConfig(cfg)
	return err
}

func (w *windows_c) StopService() error {
	_, err := w.service.Control(svc.Stop)
	return err
}

func (w *windows_c) StartService() error {
	err := w.service.Start()
	return err
}
