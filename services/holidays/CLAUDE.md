# services/holidays

Full design doc: `docs/holidays-service-checklist.md`. The simplest service in this repo — read this before over-building it.

## Hardcoded to US — no country parameter anywhere

Not config-driven, not parameterized. If multi-country ever becomes a real requirement, that's a deliberate redesign, not a flag to flip.

## API: confirm the version before implementing

`https://date.nager.at/api/v3/PublicHolidays/{year}/US` is the documented endpoint, but a newer `nagerholidays.com/api/v4/Holidays/US/{year}` also turned up during design — looks like a domain/version migration in progress. Check which is current before wiring this up; don't assume v3 is still right by the time this gets built.

## Two lifecycles, not one poll interval

- **Past years**: fetched once, ever, immutable. Never refetched — there's no reason to.
- **Current + next year**: fetched once, then re-validated on a low-frequency schedule (weekly starting point) to catch the rare legislative revision to an upcoming holiday's date.

This is the opposite shape from Weather — instead of "refresh constantly because it goes stale fast," it's "fetch once and almost never revisit." Don't apply Weather's shadow-key pattern here; it solves a problem this service doesn't have.

## No second API call for long weekends

Nager.Date has a `/LongWeekend` endpoint — don't use it. Once holiday dates are stored, "is this adjacent to a weekend" is pure day-of-week arithmetic on data already owned. A second external dependency for something computable locally isn't worth it.
