package stringutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToBytes(t *testing.T) {
	expected := map[string]map[string]float64{
		//  numeric passthrough (no suffix)
		``: map[string]float64{
			`-1`: -1,
			`0`:  0,
			`1`:  1,
			`4611686018427387903`:  4611686018427387903,
			`4611686018427387904`:  4611686018427387904,
			`4611686018427387905`:  4611686018427387905,
			`9223372036854775807`:  9223372036854775807, // beyond this overflows the positive int64 bound
			`-4611686018427387903`: -4611686018427387903,
			`-4611686018427387904`: -4611686018427387904,
			`-4611686018427387905`: -4611686018427387905,
			`-9223372036854775807`: -9223372036854775807,
			`-9223372036854775808`: -9223372036854775808, // beyond this overflows the negative int64 bound
		},

		//  suffix: b/B
		`b`: map[string]float64{
			`-1`: -1,
			`0`:  0,
			`1`:  1,
			`4611686018427387903`:  4611686018427387903,
			`4611686018427387904`:  4611686018427387904,
			`4611686018427387905`:  4611686018427387905,
			`9223372036854775807`:  9223372036854775807,
			`-4611686018427387903`: -4611686018427387903,
			`-4611686018427387904`: -4611686018427387904,
			`-4611686018427387905`: -4611686018427387905,
			`-9223372036854775807`: -9223372036854775807,
			`-9223372036854775808`: -9223372036854775808,
		},

		//  suffix: k/K
		`k`: map[string]float64{
			`-1`:               -1024,
			`0`:                0,
			`1`:                1024,
			`0.5`:              512,
			`2`:                2048,
			`9007199254740992`: 9223372036854775808,
		},

		//  suffix: m/M
		`m`: map[string]float64{
			`-1`:            -1048576,
			`0`:             0,
			`1`:             1048576,
			`0.5`:           524288,
			`8796093022208`: 9223372036854775808,
		},

		//  suffix: g/G
		`g`: map[string]float64{
			`-1`:         -1073741824,
			`0`:          0,
			`1`:          1073741824,
			`0.5`:        536870912,
			`8589934592`: 9223372036854775808,
		},

		//  suffix: t/T
		`t`: map[string]float64{
			`-1`:      -1099511627776,
			`0`:       0,
			`1`:       1099511627776,
			`0.5`:     549755813888,
			`8388608`: 9223372036854775808,
		},

		//  suffix: p/P
		`p`: map[string]float64{
			`-1`:   -1125899906842624,
			`0`:    0,
			`1`:    1125899906842624,
			`0.5`:  562949953421312,
			`8192`: 9223372036854775808,
		},

		//  suffix: e/E
		`e`: map[string]float64{
			`-1`:  -1152921504606846976,
			`0`:   0,
			`1`:   1152921504606846976,
			`0.5`: 576460752303423488,
			`8`:   9223372036854775808,
		},

		//  suffix: z/Z
		`z`: map[string]float64{
			`-1`:  -1180591620717411303424,
			`0`:   0,
			`1`:   1180591620717411303424,
			`0.5`: 590295810358705651712,
		},

		//  suffix: y/Y
		`y`: map[string]float64{
			`-1`:  -1208925819614629174706176,
			`0`:   0,
			`1`:   1208925819614629174706176,
			`0.5`: 604462909807314587353088,
		},
	}

	testExpectations := func(expectedValues map[string]float64, appendToInput string) {
		for in, out := range expectedValues {
			in = in + appendToInput

			if v, err := ToBytes(in); err == nil {
				if v != out {
					t.Errorf("Conversion error on '%s': expected %f, got %f", in, out, v)
				}
			} else {
				t.Errorf("Got error converting '%s' to bytes: %v", in, err)
			}
		}
	}

	for suffix, expectations := range expected {
		testExpectations(expectations, suffix)

		//  only unleash testing hell on higher-order conversions
		if suffix != `` && suffix != `b` {
			testExpectations(expectations, suffix+`b`)
			testExpectations(expectations, suffix+`B`)
			testExpectations(expectations, suffix+`ib`)
			testExpectations(expectations, suffix+`iB`)
		}
	}

	if v, err := ToBytes(`potato`); err == nil {
		t.Errorf("Value 'potato' inexplicably returned a value (%v)", v)
	}

	if v, err := ToBytes(`potatoG`); err == nil {
		t.Errorf("Value 'potatoG' inexplicably returned a value (%v)", v)
	}

	if v, err := ToBytes(`123X`); err == nil {
		t.Errorf("Invalid SI suffix 'X' did not error, got: %v", v)
	}
}

func TestCamelize(t *testing.T) {
	assert := require.New(t)

	tests := map[string]string{
		`Test`:        `Test`,
		`test`:        `Test`,
		`test_value`:  `TestValue`,
		`test-Value`:  `TestValue`,
		`test-Val_ue`: `TestValUe`,
		`test value`:  `TestValue`,
		`TestValue`:   `TestValue`,
		`testValue`:   `TestValue`,
		`TeSt VaLue`:  `TeStVaLue`,
	}

	for have, want := range tests {
		assert.Equal(want, Camelize(have))
	}
}

func TestUnderscore(t *testing.T) {
	assert := require.New(t)

	tests := map[string]string{
		`Test`:       `test`,
		`test`:       `test`,
		`test_value`: `test_value`,
		`test-Value`: `test_value`,
		`test value`: `test_value`,
		`TestValue`:  `test_value`,
		`testValue`:  `test_value`,
		`TeSt VaLue`: `te_st_va_lue`,
	}

	for have, want := range tests {
		assert.Equal(want, Underscore(have))
	}
}

func TestIsMixedCase(t *testing.T) {
	assert := require.New(t)

	assert.False(IsMixedCase(``))
	assert.False(IsMixedCase(`0123456789`))
	assert.False(IsMixedCase(`abcdefghijklmnopqrstuvwxyz`))
	assert.False(IsMixedCase(`abcdefghijklm0123456789nopqrstuvwxyz`))
	assert.False(IsMixedCase(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`))
	assert.False(IsMixedCase(`ABCDEFGHIJKLM0123456789NOPQRSTUVWXYZ`))
	assert.False(IsMixedCase(` ABCDEFGHIJKLM 0123456789 NOPQRSTUVWXYZ `))
	assert.False(IsMixedCase(`сою́з`))
	assert.False(IsMixedCase(`СОЮ́З`))

	assert.True(IsMixedCase(`AbCdEfGhIjKlMnOpQrStUvWxYz`))
	assert.True(IsMixedCase(`ABCDEFGHIJKLM0123456789nopqrstuvwxyz`))
	assert.True(IsMixedCase(`Сою́з`))
}

func TestIsHexadecimal(t *testing.T) {
	assert := require.New(t)

	for i := 0; i < 16; i++ {
		assert.True(IsHexadecimal(fmt.Sprintf("%x", i), -1))
		assert.True(IsHexadecimal(fmt.Sprintf("%X", i), -1))
	}

	for i := 10; i < 16; i++ {
		assert.False(IsHexadecimal(fmt.Sprintf("%x%X", i, i), -1))
		assert.False(IsHexadecimal(fmt.Sprintf("%X%x", i, i), -1))
		assert.False(IsHexadecimal(fmt.Sprintf("%x", i), 2))
		assert.False(IsHexadecimal(fmt.Sprintf("%X", i), 2))
	}

	assert.True(IsHexadecimal(`abc123`, -1))
	assert.True(IsHexadecimal(`ABC123`, -1))
	assert.True(IsHexadecimal(`abc123`, 6))
	assert.True(IsHexadecimal(`ABC123`, 6))
	assert.True(IsHexadecimal(`b26252862a11dd3221427bdbae6025604b1760e4`, 40))

	assert.False(IsHexadecimal(`aBc123`, -1))
	assert.False(IsHexadecimal(`abc123`, 32))
	assert.False(IsHexadecimal(`ABC123`, 32))

}

func TestThousandify(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, Thousandify(``, ``, ``))
	assert.Equal(`0`, Thousandify(`0`, ``, ``))
	assert.Equal(`1`, Thousandify(`1`, ``, ``))
	assert.Equal(`1,000`, Thousandify(`1000`, ``, ``))
	assert.Equal(`1,000.000`, Thousandify(`1000.000`, ``, ``))
	assert.Equal(`1,000.001`, Thousandify(`1000.001`, ``, ``))
	assert.Equal(`9,223,372,036,854,775,807`, Thousandify(`9223372036854775807`, ``, ``))
	assert.Equal(`-9,223,372,036,854,775,809`, Thousandify(`-9223372036854775809`, ``, ``))
	assert.Equal(
		`-9,223,372,036,854,775,809,922,337,203,685,477,580,992,233,720,368,547,758,099,223,372,036,854,775,809`,
		Thousandify(`-9223372036854775809922337203685477580992233720368547758099223372036854775809`, ``, ``),
	)

	assert.Equal(`0`, Thousandify(0, ``, ``))
	assert.Equal(`1`, Thousandify(1, ``, ``))
	assert.Equal(`1,000`, Thousandify(1000, ``, ``))
	assert.Equal(`1,000`, Thousandify(1000.000, ``, ``))
	assert.Equal(`1,000.001`, Thousandify(1000.001, ``, ``))

	assert.Equal(`9,223,372,036,854,775,807`, Thousandify(9223372036854775807, ``, ``))
	assert.Equal(`-9,223,372,036,854,775,808`, Thousandify(-9223372036854775808, ``, ``))
}

func TestLongestCommonPrefix(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, LongestCommonPrefix(nil))
	assert.Equal(`interstellar`, LongestCommonPrefix([]string{
		`interstellar`,
	}))

	assert.Equal(`inters`, LongestCommonPrefix([]string{
		`interstellar`,
		`interspace`,
		`interstitial`,
	}))

	assert.Equal(`inter`, LongestCommonPrefix([]string{
		`interstellar`,
		`interspace`,
		`interstitial`,
		`interesting`,
		`interest`,
	}))

	assert.Equal(`test.`, LongestCommonPrefix([]string{
		`test.value`,
		`test.debug`,
		`test.test`,
	}))
}

func TestRelaxedEqual(t *testing.T) {
	assert := require.New(t)

	eq, err := RelaxedEqual(nil, nil)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(1, 1)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(int(1), int64(1))
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(float64(1), byte(1))
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(float64(1.00), `1`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(true, true)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(false, false)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`true`, `on`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`true`, `yes`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(`boo`, `boo`)
	assert.NoError(err)
	assert.True(eq)

	eq, err = RelaxedEqual(1, true)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(true, false)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(false, true)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`true`, `no`)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`false`, `yes`)
	assert.NoError(err)
	assert.False(eq)

	eq, err = RelaxedEqual(`boo`, `Boo`)
	assert.NoError(err)
	assert.False(eq)
}

func TestSplitWords(t *testing.T) {
	assert := require.New(t)

	assert.Equal([]string{
		`Goldenrod-adorned`,
		`log`,
		`.`,
	}, SplitWords(`Goldenrod-adorned log.`))

	assert.Equal([]string{
		`Goldenrod`,
		`adorned`,
		`log`,
		`.`,
	}, SplitWords(`Goldenrod adorned log.`))
}


func TestElideWords(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, ElideWords(``, 0))
	assert.Equal(`.`, ElideWords(`.`, 1))
	assert.Equal(`test`, ElideWords(`test.`, 1))

	assert.Equal(
		`This is the song that never ends`,
		ElideWords(
			`This is the song that never ends, it just goes on and on my friends.`,
			7,
		),
	)
}
