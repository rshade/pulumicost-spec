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
 * Get all valid UsageProfile enum values.
 *
 * @returns An array containing every known `UsageProfile` value
 */
export function getAllUsageProfiles(): ReadonlyArray<UsageProfile> {
  return allUsageProfiles;
}

/**
 * Determines whether the provided value is one of the known UsageProfile enum values.
 *
 * @returns `true` if the profile is `UNSPECIFIED`, `PROD`, `DEV`, or `BURST`, `false` otherwise.
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
 * Parse a string into a UsageProfile enum value.
 *
 * Accepts case-insensitive inputs and common variants: the empty string maps to `UNSPECIFIED`; `"dev"` or `"development"` map to `DEV`; `"prod"` or `"production"` map to `PROD`; `"burst"` maps to `BURST`.
 *
 * @param s - The input string to parse.
 * @returns The corresponding `UsageProfile` enum value.
 * @throws Error if `s` is not a recognized usage profile (error message: `unknown usage profile: "<input>"`).
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
 * Get the lowercase display string for a usage profile.
 *
 * @returns `"unspecified"`, `"prod"`, `"dev"`, or `"burst"` for known profiles; otherwise `unknown(<profile>)`
 */
export function usageProfileString(profile: UsageProfile): string {
  return stringMap.get(profile) ?? `unknown(${profile})`;
}

/**
 * Normalize a UsageProfile to a known enum value.
 *
 * Returns `profile` when it is one of the recognized UsageProfile values; otherwise returns `UsageProfile.UNSPECIFIED`, enabling forward compatibility with newer spec values.
 *
 * @returns `profile` if it is a known UsageProfile; `UsageProfile.UNSPECIFIED` otherwise.
 */
export function normalizeUsageProfile(profile: UsageProfile): UsageProfile {
  if (isValidUsageProfile(profile)) {
    return profile;
  }
  console.warn(`Unknown usage profile (${profile}), treating as UNSPECIFIED`);
  return UsageProfile.UNSPECIFIED;
}

/**
 * Get the default monthly usage hours for a usage profile.
 *
 * Mapping:
 * - PROD => 730
 * - DEV => 160
 * - BURST => 200
 * - UNSPECIFIED => 730 (defaults to production)
 *
 * @param profile - The usage profile to evaluate
 * @returns The default number of hours per month for the specified profile
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