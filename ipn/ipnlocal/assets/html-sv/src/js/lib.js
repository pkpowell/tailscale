const Base2 = {
    factor: 1024, 
    suffix: "iB"
}
const Base10 = {
    factor: 1000, 
    suffix: "B"
}
const PRX = {
    0: {
        short: 'K',
        full:  "kilo",
    },
    1: {
        short: 'M',
        full:  "mega",
    },
    2: {
        short: 'G',
        full:  "giga",
    },
    3: {
        short: 'T',
        full:  "tera",
    },
    4: {
        short: 'P',
        full:  "peta",
    },
    5: {
        short: 'E',
        full:  "exa",
    },
    // these are just for show. int64 only reaches 8EB
    6: {
        short: 'Y',
        full:  "yotta",
    },
    7: {
        short: 'Z',
        full:  "zetta",
    },
}

// FormatBytes converts bytes to KB, MiB etc without Math lib
const FormatBytes = (b, u) => {
	// bytes only
	if (b === 0) {
		return "–"
	}

	if (b < u.factor) {
		return b + "B"
	}

	let div = u.factor
	let exp = 0

	for (let n = b / u.factor; n >= u.factor; n /= u.factor) {
		// grow the divisor
		div *= u.factor
		exp++
	}
	return (b / div).toFixed(2), PRX[exp].short, u.suffix
}

export  {
    FormatBytes
}