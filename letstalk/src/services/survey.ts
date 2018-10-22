import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { SURVEY_ROUTE } from "./constants";
import {Survey} from "../models/survey";

export interface SurveyService {
  postSurveyResponses(surveyWithResponses: Survey): Promise<Object>;
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

  async postSurveyResponses(surveyWithResponses: Survey): Promise<Object> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.post(SURVEY_ROUTE, surveyWithResponses, sessionToken);
  }
}

const surveyService = new RemoteSurveyService(requestor, auth);
export default surveyService;
