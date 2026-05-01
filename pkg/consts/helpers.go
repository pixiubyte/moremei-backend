package consts

func IsStatus(s int64, target int64) bool {
    return s == target
}

func IsContainStatus(s int64, targets []int64) bool {
    if len(targets) == 0 {
        return false
    }
    for _, target := range targets {
        if IsStatus(s, target) {
            return true
        }
    }

    return false
}

func IsStatusDaft(s int64) bool {
    return IsStatus(s, StatusDaft)
}

func IsStatusNew(s int64) bool {
    return IsStatus(s, StatusNew)
}

func IsStatusProcessing(s int64) bool {
    return IsStatus(s, StatusProcessing)
}

func IsStatusSuccessful(s int64) bool {
    return IsStatus(s, StatusSuccessful)
}

func IsStatusFailure(s int64) bool {
    return IsStatus(s, StatusFailure)
}

func IsFlag(f uint64, target uint64) bool {
    return f == target
}

func IsFlagTrue(f uint64) bool {
    return IsFlag(f, FlagTrue)
}

func IsFlagFalse(f uint64) bool {
    return IsFlag(f, FlagFalse)
}
