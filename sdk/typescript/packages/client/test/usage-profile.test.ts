import { describe, it, expect } from 'vitest';
import { UsageProfile } from '../src/generated/finfocus/v1/enums_pb.js';
import {
  getAllUsageProfiles,
  isValidUsageProfile,
  parseUsageProfile,
  usageProfileString,
  normalizeUsageProfile,
  defaultMonthlyHours,
} from '../src/utils/usage-profile.js';

describe('UsageProfile helpers', () => {
  describe('getAllUsageProfiles', () => {
    it('returns all four profile values', () => {
      const profiles = getAllUsageProfiles();
      expect(profiles).toHaveLength(4);
      expect(profiles).toContain(UsageProfile.UNSPECIFIED);
      expect(profiles).toContain(UsageProfile.PROD);
      expect(profiles).toContain(UsageProfile.DEV);
      expect(profiles).toContain(UsageProfile.BURST);
    });
  });

  describe('isValidUsageProfile', () => {
    it.each([
      [UsageProfile.UNSPECIFIED, true],
      [UsageProfile.PROD, true],
      [UsageProfile.DEV, true],
      [UsageProfile.BURST, true],
      [999 as UsageProfile, false],
      [-1 as UsageProfile, false],
    ])('isValidUsageProfile(%s) returns %s', (profile, expected) => {
      expect(isValidUsageProfile(profile)).toBe(expected);
    });
  });

  describe('parseUsageProfile', () => {
    it.each([
      ['', UsageProfile.UNSPECIFIED],
      ['unspecified', UsageProfile.UNSPECIFIED],
      ['dev', UsageProfile.DEV],
      ['DEV', UsageProfile.DEV],
      ['Dev', UsageProfile.DEV],
      ['development', UsageProfile.DEV],
      ['prod', UsageProfile.PROD],
      ['PROD', UsageProfile.PROD],
      ['production', UsageProfile.PROD],
      ['burst', UsageProfile.BURST],
      ['BURST', UsageProfile.BURST],
    ])('parseUsageProfile("%s") returns %s', (input, expected) => {
      expect(parseUsageProfile(input)).toBe(expected);
    });

    it('throws for unknown strings', () => {
      expect(() => parseUsageProfile('invalid')).toThrow('unknown usage profile: "invalid"');
      expect(() => parseUsageProfile('test')).toThrow('unknown usage profile');
    });
  });

  describe('usageProfileString', () => {
    it.each([
      [UsageProfile.UNSPECIFIED, 'unspecified'],
      [UsageProfile.PROD, 'prod'],
      [UsageProfile.DEV, 'dev'],
      [UsageProfile.BURST, 'burst'],
    ])('usageProfileString(%s) returns "%s"', (profile, expected) => {
      expect(usageProfileString(profile)).toBe(expected);
    });

    it('returns unknown format for unrecognized values', () => {
      expect(usageProfileString(999 as UsageProfile)).toBe('unknown(999)');
    });
  });

  describe('normalizeUsageProfile', () => {
    it('passes through known profiles unchanged', () => {
      expect(normalizeUsageProfile(UsageProfile.PROD)).toBe(UsageProfile.PROD);
      expect(normalizeUsageProfile(UsageProfile.DEV)).toBe(UsageProfile.DEV);
      expect(normalizeUsageProfile(UsageProfile.BURST)).toBe(UsageProfile.BURST);
      expect(normalizeUsageProfile(UsageProfile.UNSPECIFIED)).toBe(UsageProfile.UNSPECIFIED);
    });

    it('normalizes unknown values to UNSPECIFIED', () => {
      expect(normalizeUsageProfile(999 as UsageProfile)).toBe(UsageProfile.UNSPECIFIED);
      expect(normalizeUsageProfile(-1 as UsageProfile)).toBe(UsageProfile.UNSPECIFIED);
    });
  });

  describe('defaultMonthlyHours', () => {
    it.each([
      [UsageProfile.PROD, 730],
      [UsageProfile.DEV, 160],
      [UsageProfile.BURST, 200],
      [UsageProfile.UNSPECIFIED, 730],
    ])('defaultMonthlyHours(%s) returns %d', (profile, expected) => {
      expect(defaultMonthlyHours(profile)).toBe(expected);
    });

    it('returns production hours for unknown profiles', () => {
      expect(defaultMonthlyHours(999 as UsageProfile)).toBe(730);
    });
  });
});
