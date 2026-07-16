package main

import (
	"fmt"
	"log"
)

func executeRemoteCommand(cmd RemoteCommand) RemoteResult {
	log.Printf("agent executing op=%s account=%s id=%s", cmd.Op, cmd.Account, cmd.ID)
	result := RemoteResult{ID: cmd.ID, OK: false}

	switch cmd.Op {
	case "get_state":
		locked, lockErr := getLockstate(cmd.Account)
		minutes, timeErr := getLogintime(cmd.Account)
		remaining, timeoutErr := getTimeout(cmd.Account)

		result.LockState = &locked
		dur := minutes
		result.Duration = &dur
		rem := remaining
		result.Remaining = &rem

		if lockErr != nil {
			result.Message = "Error getting lockstate"
			result.Error = lockErr.Error()
			break
		}
		if timeErr != nil {
			result.Message = "Error getting logintime"
			result.Error = timeErr.Error()
			break
		}
		if timeoutErr != nil {
			result.Message = "Error getting timeout"
			result.Error = timeoutErr.Error()
			break
		}

		result.OK = true
		if minutes == -1 {
			result.Message = "Account not logged in"
		} else {
			result.Message = fmt.Sprintf("Account logged in for %d min", minutes)
		}
	case "get_lock":
		locked, err := getLockstate(cmd.Account)
		if err != nil {
			result.Message = "Error getting lockstate"
			result.Error = err.Error()
			break
		}
		result.OK = true
		result.LockState = &locked
		if locked {
			result.Message = "Account locked"
		} else {
			result.Message = "Account unlocked"
		}
	case "set_lock":
		if cmd.LockState == nil {
			result.Message = "Missing lockstate"
			result.Error = "lockstate missing"
			break
		}
		err := setLockstate(cmd.Account, *cmd.LockState)
		if err != nil {
			result.Message = "Error setting account lockstate"
			result.Error = err.Error()
			break
		}
		result.OK = true
		if *cmd.LockState {
			result.Message = "Account locked"
		} else {
			result.Message = "Account unlocked"
		}
		result.LockState = cmd.LockState
	case "get_time":
		minutes, err := getLogintime(cmd.Account)
		dur := minutes
		result.Duration = &dur
		if err != nil {
			result.Message = "Error getting logintime"
			result.Error = err.Error()
			break
		}
		result.OK = true
		if minutes == -1 {
			result.Message = "Account not logged in"
		} else {
			result.Message = fmt.Sprintf("Account logged in for %d min", minutes)
		}
	case "killall":
		err := killall(cmd.Account)
		if err != nil {
			result.Message = "Error killing all processes"
			result.Error = err.Error()
			break
		}
		result.OK = true
		result.Message = "All processes killed"
	case "get_timeout":
		remaining, err := getTimeout(cmd.Account)
		rem := remaining
		result.Remaining = &rem
		if err != nil {
			result.Message = "Timeout unreadable"
			result.Error = err.Error()
			break
		}
		result.OK = true
		if remaining == -1 {
			result.Message = "Timeout not found"
		} else {
			result.Message = fmt.Sprintf("Timeout is %d", remaining)
		}
	case "set_timeout":
		err := setTimeout(cmd.Duration, cmd.Account)
		if err != nil {
			result.Message = "Error setting timeout"
			result.Error = err.Error()
			break
		}
		result.OK = true
		rem := cmd.Duration
		result.Remaining = &rem
		result.Message = fmt.Sprintf("Timeout is %d", cmd.Duration)
	case "del_timeout":
		success, err := delTimeout(cmd.Account)
		if err != nil {
			result.Message = "Timeout not deleted"
			result.Error = err.Error()
			break
		}
		result.OK = true
		if success {
			result.Message = "Timeout unset"
		} else {
			result.Message = "Timeout not found"
		}
	default:
		result.Message = "Unknown operation"
		result.Error = "unknown op"
	}

	log.Printf("agent finished op=%s account=%s id=%s ok=%t message=%q", cmd.Op, cmd.Account, cmd.ID, result.OK, result.Message)
	return result
}
