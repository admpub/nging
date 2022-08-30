// +build darwin

package system

/*
const _Csconst_AMBIENT_AIR_0 = "TA0P"
const _Csconst_AMBIENT_AIR_1 = "TA1P"
const _Csconst_CPU_0_DIODE = "TC0D"
const _Csconst_CPU_0_HEATSINK = "TC0H"
const _Csconst_CPU_0_PROXIMITY = "TC0P"
const _Csconst_ENCLOSURE_BASE_0 = "TB0T"
const _Csconst_ENCLOSURE_BASE_1 = "TB1T"
const _Csconst_ENCLOSURE_BASE_2 = "TB2T"
const _Csconst_ENCLOSURE_BASE_3 = "TB3T"
const _Csconst_GPU_0_DIODE = "TG0D"
const _Csconst_GPU_0_HEATSINK = "TG0H"
const _Csconst_GPU_0_PROXIMITY = "TG0P"
const _Csconst_HARD_DRIVE_BAY = "TH0P"
const _Csconst_MEMORY_SLOTS_PROXIMITY = "TM0P"
const _Csconst_MEMORY_SLOT_0 = "TM0S"
const _Csconst_NORTHBRIDGE = "TN0H"
const _Csconst_NORTHBRIDGE_DIODE = "TN0D"
const _Csconst_NORTHBRIDGE_PROXIMITY = "TN0P"
const _Csconst_THUNDERBOLT_0 = "TI0P"
const _Csconst_THUNDERBOLT_1 = "TI1P"
const _Csconst_WIRELESS_MODULE = "TW0P"
*/

var temperatureKeys = map[string]string{
	`TA0P`: "ambient_air_0",
	`TA1P`: "ambient_air_1",
	`TC0D`: "cpu_0_diode",
	`TC0H`: "cpu_0_heatsink",
	`TC0P`: "cpu_0_proximity",
	`TB0T`: "enclosure_base_0",
	`TB1T`: "enclosure_base_1",
	`TB2T`: "enclosure_base_2",
	`TB3T`: "enclosure_base_3",
	`TG0D`: "gpu_0_diode",
	`TG0H`: "gpu_0_heatsink",
	`TG0P`: "gpu_0_proximity",
	`TH0P`: "hard_drive_bay",
	`TM0S`: "memory_slot_0",
	`TM0P`: "memory_slots_proximity",
	`TN0H`: "northbridge",
	`TN0D`: "northbridge_diode",
	`TN0P`: "northbridge_proximity",
	`TI0P`: "thunderbolt_0",
	`TI1P`: "thunderbolt_1",
	`TW0P`: "wireless_module",
}

func init() {
	keyFormatter = func(s string) string {
		key, ok := temperatureKeys[s]
		if !ok {
			return s
		}
		return key
	}
}
