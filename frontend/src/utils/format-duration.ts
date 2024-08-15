const divMod = (n: number, m: number) => [Math.floor(n / m), n % m];

export function createDurationFormatter(
  unitDisplay: Intl.NumberFormatOptions["unitDisplay"] = "long",
  locale: Intl.LocalesArgument = navigator.language,
  showMs = false,
): (ms: number) => string {
  const timeUnitFormatter = (
      locale: Intl.LocalesArgument,
      unit: Intl.NumberFormatOptions["unit"],
      unitDisplay: Intl.NumberFormatOptions["unitDisplay"],
    ) => Intl.NumberFormat(locale, { style: "unit", unit, unitDisplay }).format,
    fmtDays = timeUnitFormatter(locale, "day", unitDisplay),
    fmtHours = timeUnitFormatter(locale, "hour", unitDisplay),
    fmtMinutes = timeUnitFormatter(locale, "minute", unitDisplay),
    fmtSeconds = timeUnitFormatter(locale, "second", unitDisplay),
    fmtMilliseconds = timeUnitFormatter(locale, "millisecond", unitDisplay),
    fmtList = new Intl.ListFormat(locale, {
      style: "long",
      type: "conjunction",
    });
  return (milliseconds: number) => {
    let days, hours, minutes, seconds;
    [days, milliseconds] = divMod(milliseconds, 864e5);
    [hours, milliseconds] = divMod(milliseconds, 36e5);
    [minutes, milliseconds] = divMod(milliseconds, 6e4);
    [seconds, milliseconds] = divMod(milliseconds, 1e3);

    const parts = [
      days ? fmtDays(days) : null,
      hours ? fmtHours(hours) : null,
      minutes ? fmtMinutes(minutes) : null,
      seconds ? fmtSeconds(seconds) : null,
      milliseconds && showMs ? fmtMilliseconds(milliseconds) : null,
    ];

    if (unitDisplay == "narrow") {
      return parts.join("");
    }

    return fmtList.format(parts.filter((v) => v !== null));
  };
}
