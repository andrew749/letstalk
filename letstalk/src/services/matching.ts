import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { MatchingData } from '../models/matching';
import { MATCHING_ROUTE } from "./constants";

export interface MatchingService {
  getMatchingRequest(matchedUserId: int): Promise<MatchingData>;
  putMatchingRequest(matchingData: MatchingData): Promise<MatchingData>;
}

export class RemoteMatchingService implements MatchingService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async getMatchingRequest(matchedUserId: int): Promise<MatchingData> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.get(`${MATCHING_ROUTE}/${matchedUserId}`, sessionToken);
  }

  async putMatchingRequest(matchingData: MatchingData): Promise<MatchingData> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.put(MATCHING_ROUTE, matchingData, sessionToken);
  }
}

const matchingService = new RemoteMatchingService(requestor, auth);
export default matchingService;
