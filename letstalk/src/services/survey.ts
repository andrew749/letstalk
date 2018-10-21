import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import {SURVEY_RESPONSES_ROUTE, SURVEY_ROUTE} from "./constants";
import {Survey, SurveyResponses} from "../models/survey";

export interface SurveyService {
  postSurveyResponses(responses: SurveyResponses): Promise<Object>;
}

export class RemoteSurveyService implements SurveyService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async getSurvey(): Promise<Survey> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.post(SURVEY_ROUTE, responses, sessionToken);
  }

  async postSurveyResponses(responses: SurveyResponses): Promise<Object> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.post(SURVEY_RESPONSES_ROUTE, responses, sessionToken);
  }
}

const surveyService = new RemoteSurveyService(requestor, auth);
export default surveyService;
