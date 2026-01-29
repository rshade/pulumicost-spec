import { UsageProfile } from "../generated/finfocus/v1/enums_pb.js";

/** All valid UsageProfile enum values. */
const allUsageProfiles: ReadonlyArray<UsageProfile> = [
  UsageProfile.UNSPECIFIED,
  UsageProfile.PROD,
  UsageProfile.DEV,
  UsageProfile.BURST,
];

/** Default monthly hours per profile. */
const HOURS_PROD = 730;
const HOURS_DEV = 160;
const HOURS_BURST = 200;

/**
 * Returns all valid UsageProfile enum values.
 */
export function getAllUsageProfiles(): ReadonlyArray<UsageProfile> {
  return allUsageProfiles;
}

/**
 * Checks if the given value is a known UsageProfile enum value.
 * Returns true for UNSPECIFIED, PROD, DEV, and BURST.
 */
export function isValidUsageProfile(profile: UsageProfile): boolean {
  return allUsageProfiles.includes(profile);
}

/** Map from lowercase string to UsageProfile for parsing. */
const parseMap: ReadonlyMap<string, UsageProfile> = new Map([
  ["unspecified", UsageProfile.UNSPECIFIED],
  ["prod", UsageProfile.PROD],
  ["production", UsageProfile.PROD],
  ["dev", UsageProfile.DEV],
  ["development", UsageProfile.DEV],
  ["burst", UsageProfile.BURST],
]);

/**
 * Parses a string into a UsageProfile enum value.
 * Supports case-insensitive matching and common variants:
 *   - "dev", "development" → DEV
 *   - "prod", "production" → PROD
 *   - "burst" → BURST
 *   - "unspecified", "" → UNSPECIFIED
 *
 * @throws Error for unrecognized strings.
 */
export function parseUsageProfile(s: string): UsageProfile {
  if (s === "") {
    return UsageProfile.UNSPECIFIED;
  }

  const profile = parseMap.get(s.toLowerCase());
  if (profile !== undefined) {
    return profile;
  }

  throw new Error(`unknown usage profile: "${s}"`);
}

/** Map from UsageProfile to lowercase string for display. */
const stringMap: ReadonlyMap<UsageProfile, string> = new Map([
  [UsageProfile.UNSPECIFIED, "unspecified"],
  [UsageProfile.PROD, "prod"],
  [UsageProfile.DEV, "dev"],
  [UsageProfile.BURST, "burst"],
]);

/**
 * Returns a lowercase string representation of the profile.
 * For known profiles: "unspecified", "prod", "dev", "burst".
 * For unknown profiles: "unknown(<value>)".
 */
export function usageProfileString(profile: UsageProfile): string {
  return stringMap.get(profile) ?? `unknown(${profile})`;
}

/**
 * Returns the profile if it's a known value, or UNSPECIFIED for unknown values.
 * Enables forward compatibility when receiving profile values from newer spec versions.
 */
export function normalizeUsageProfile(profile: UsageProfile): UsageProfile {
  if (isValidUsageProfile(profile)) {
    return profile;
  }
  return UsageProfile.UNSPECIFIED;
}

/**
 * Returns the default monthly usage hours for a profile.
 *   - PROD: 730 hours (24/7 operation)
 *   - DEV: 160 hours (~8 hours/day, 5 days/week)
 *   - BURST: 200 hours (plugin discretion)
 *   - UNSPECIFIED: 730 hours (defaults to production)
 *
 * Plugins have discretion to use different values based on their resource types.
 */
export function defaultMonthlyHours(profile: UsageProfile): number {
  switch (profile) {
    case UsageProfile.PROD:
      return HOURS_PROD;
    case UsageProfile.DEV:
      return HOURS_DEV;
    case UsageProfile.BURST:
      return HOURS_BURST;
    case UsageProfile.UNSPECIFIED:
      return HOURS_PROD;
    default:
      return HOURS_PROD;
  }
}
