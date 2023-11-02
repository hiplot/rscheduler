package scheduler

import (
	"fmt"
	"rscheduler/processor"
)

func FindProcByProcID(id string) *processor.Proc {
	RScheduler.Lock.RLock()
	defer RScheduler.Lock.RUnlock()

	for _, procList := range RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if proc.ID == id {
				return proc
			}
		}
	}

	return nil
}

func FindProcByTaskID(id string) *processor.Proc {
	RScheduler.Lock.RLock()
	defer RScheduler.Lock.RUnlock()

	for _, procList := range RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if proc.Task != nil && proc.Task.ID == id {
				return proc
			}
		}
	}

	return nil
}

func KillProcByProcID(id string, force bool) error {
	RScheduler.Lock.Lock()
	defer RScheduler.Lock.Unlock()

	for _, procList := range RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if proc.ID == id {
				if force || proc.IsIdle() {
					if err := proc.ForceClose(); err != nil {
						return err
					}
					procList.Remove(i)
					return nil
				}
				proc.SetPreDelete()
				return nil
			}
		}
	}

	return fmt.Errorf("can not find proc by proc id: %s", id)
}

func KillProcByTaskID(id string) error {
	RScheduler.Lock.Lock()
	defer RScheduler.Lock.Unlock()

	for _, procList := range RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if proc.Task != nil && proc.Task.ID == id {
				if err := proc.ForceClose(); err != nil {
					return err
				}
				procList.Remove(i)
				return nil
			}
		}
	}

	return fmt.Errorf("can not find proc by task id: %s", id)
}
