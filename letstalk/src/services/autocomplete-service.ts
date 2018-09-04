import requestor, { Requestor } from './requests';
import { SimpleTrait } from '../models/simple-trait';
import { MultiTrait } from '../models/multi-trait';
import { Role, Organization } from '../models/position';
import {
  AUTOCOMPLETE_SIMPLE_TRAIT_ROUTE,
  AUTOCOMPLETE_ROLE_ROUTE,
  AUTOCOMPLETE_ORGANIZATION_ROUTE,
  AUTOCOMPLETE_MULTI_TRAIT_ROUTE,
} from './constants';

interface AutocompleteRequest {
  readonly prefix: string;
  readonly size: number;
}

class AutocompleteService {
  private requestor: Requestor

  constructor(requestor: Requestor) {
    this.requestor = requestor;
  }

  private async doRequest(url: string, prefix: string, size: number): Promise<any> {
    const req: AutocompleteRequest = { prefix, size };
    const res = await this.requestor.post(url, req);
    return res
  }

  async autocompleteSimpleTrait(prefix: string, size: number): Promise<Array<SimpleTrait>> {
    return this.doRequest(AUTOCOMPLETE_SIMPLE_TRAIT_ROUTE, prefix, size);
  }

  async autocompleteOrganization(prefix: string, size: number): Promise<Array<Organization>> {
    return this.doRequest(AUTOCOMPLETE_ORGANIZATION_ROUTE, prefix, size);
  }

  async autocompleteRole(prefix: string, size: number): Promise<Array<Role>> {
    return this.doRequest(AUTOCOMPLETE_ROLE_ROUTE, prefix, size);
  }

  async autocompleteMultiTrait(prefix: string, size: number): Promise<Array<MultiTrait>> {
    return this.doRequest(AUTOCOMPLETE_MULTI_TRAIT_ROUTE, prefix, size);
  }
}

const autocompleteService = new AutocompleteService(requestor);
export default autocompleteService;
