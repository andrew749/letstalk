export enum MatchingStateId {
    Unverified = 0,
    Verified,
    Expired,
}

export function matchingStateIdToString(matchingStateId: MatchingStateId): string {
  switch (matchingStateId) {
    case MatchingStateId.Unverified:
      return 'unverified';
  case MatchingStateId.Verified:
      return 'verified';
  case MatchingStateId.Expired:
      return 'expired';
    default:
      const _: never = matchingStateId;
  }
}

export interface MatchingData {
    readonly Mentor: number,
    readonly Mentee: number,
    Secret: string,
    State: MatchingStateId
}
