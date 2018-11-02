import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { SURVEY_ROUTE } from "./constants";
import {Survey} from "../models/survey";

export interface SurveyService {
  postSurveyResponses(surveyWithResponses: Survey): Promise<Object>;
}

export const GROUP_GENERIC = "generic";

export class RemoteSurveyService implements SurveyService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async getSurvey(group: string): Promise<Survey> {
    const sessionToken = await auth.getSessionToken();
    const surveyGroupRoute = `${SURVEY_ROUTE}/${group}`;
    return await this.requestor.get(surveyGroupRoute, sessionToken);
  }

  async postSurveyResponses(surveyWithResponses: Survey): Promise<Object> {
    const sessionToken = await auth.getSessionToken();
    return await this.requestor.post(SURVEY_ROUTE, surveyWithResponses, sessionToken);
  }
}

const surveyService = new RemoteSurveyService(requestor, auth);
export default surveyService;
