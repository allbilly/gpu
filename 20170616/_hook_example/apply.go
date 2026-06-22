mapWGLog := timing.NewMapWGLog(logger)
computeUnit.Scheduler.AcceptHook(mapWGLog)