package status

import "testing"

func Test_getMajorityNumber(t *testing.T) {

	// act
	numbersAndExpectedResults := []struct {
		number         int
		expectedResult int
	}{{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{5, 3},
		{6, 4},
		{7, 4},
		{8, 5},
		{9, 5},
		{10, 6},
		{11, 6},
		{12, 7},
		{13, 7},
	}

	// arrange
	for _, entry := range numbersAndExpectedResults {

		// act
		result := getMajorityNumber(entry.number)

		// assert
		if result != entry.expectedResult {
			t.Fail()
			t.Logf("getMajorityNumber(%d) should have returned %q but returned %d instead.", entry.number, entry.expectedResult, result)
		}
	}

}

func Test_getMajorityStatus_AllStatusesAreNotSet_NotSetStatusIsReturned(t *testing.T) {
	// arrange
	statuses := []Status{
		NotSet,
		NotSet,
		NotSet,
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result != NotSet {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", NotSet, result)
	}
}

func Test_getMajorityStatus_MajorityAvailable_MajorIsReturned(t *testing.T) {
	// arrange
	statuses := []Status{
		NotInstalled,
		NotInstalled,
		NotInstalled,
		NotSet,
		NotSet,
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if !exists || result != NotInstalled {
		t.Fail()
		t.Logf("getMajorityStatus should have returned %q but returned %q instead.", NotInstalled, result)
	}
}

func Test_getMajorityStatus_NoMajority_ResultIsFalse(t *testing.T) {
	// arrange
	statuses := []Status{
		NotSet,
		Outdated,
		UpToDate,
	}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if exists || result != NotSet {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though there is no majority.")
	}
}

func Test_getMajorityStatus_NoStatuses_ResultIsFalse(t *testing.T) {
	// arrange
	statuses := []Status{}

	// act
	exists, result := getMajorityStatus(statuses)

	// assert
	if exists || result != NotSet {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though no statuses were given.")
	}
}

func Test_getMajorityStatus_Nil_ResultIsFalse(t *testing.T) {
	// act
	exists, result := getMajorityStatus(nil)

	// assert
	if exists || result != NotSet {
		t.Fail()
		t.Logf("getMajorityStatus returned a status even though no statuses were given.")
	}
}

func Test_CalculateGlobalStatus_EmptyList_ResultStatusIsNotSet(t *testing.T) {
	// arrange
	statuses := []Status{}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result != NotSet {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, NotSet)
	}
}

func Test_CalculateGlobalStatus_MixedStatuses_ResultStatusIsNotSet(t *testing.T) {
	// arrange
	statuses := []Status{
		NotSet,
		NotSet,
		UpToDate,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result != NotSet {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, NotSet)
	}
}

func Test_CalculateGlobalStatus_AllStatusesAreUpToDate_ResultStatusIsUpToDate(t *testing.T) {
	// arrange
	statuses := []Status{
		UpToDate,
		UpToDate,
		UpToDate,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result != UpToDate {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, UpToDate)
	}
}

func Test_CalculateGlobalStatus_AllStatusesAreOutdated_ResultStatusIsOutdated(t *testing.T) {
	// arrange
	statuses := []Status{
		Outdated,
		Outdated,
		Outdated,
	}

	// act
	result := CalculateGlobalStatus(statuses)

	// assert
	if result != Outdated {
		t.Fail()
		t.Logf("CalculateGlobalStatus returned %s instead of %s", result, Outdated)
	}
}
