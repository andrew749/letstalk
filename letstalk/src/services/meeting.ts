import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { MeetingConfirmation } from '../models/meeting';
import { MEETING_CONFIRMATION_ROUTE } from "./constants";

export interface MeetingService {
  postMeetingConfirmation(meetingConfirmation: MeetingConfirmation): Promise<MeetingConfirmation>;
}

export class RemoteMeetingService implements MeetingService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async postMeetingConfirmation(meetingConfirmation: MeetingConfirmation): Promise<MeetingConfirmation> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.post(MEETING_CONFIRMATION_ROUTE, meetingConfirmation, sessionToken);
  }
}

const meetingService = new RemoteMeetingService(requestor, auth);
export default meetingService;
