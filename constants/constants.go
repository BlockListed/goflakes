package constants

const (
	IdLength                = 63
	TimestampLength         = 41
	InstanceLength          = 10
	SequenceLength          = 12
	TimestapShift           = IdLength - TimestampLength
	InstanceShift           = TimestapShift - InstanceLength
	LatestStorableTime      = (1 << TimestampLength) - 1
	BiggestStorableInstance = (1 << InstanceLength) - 1
	BiggestStorableSequence = (1 << SequenceLength) - 1
	ResetSequence           = 1 << SequenceLength
	TimestampMask           = LatestStorableTime << TimestapShift
	InstanceMask            = BiggestStorableInstance << InstanceShift
	SequenceMask            = BiggestStorableSequence
)
