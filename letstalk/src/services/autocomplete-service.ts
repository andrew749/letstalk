import requestor, { Requestor } from './requests';
import { AUTOCOMPLETE_SIMPLE_TRAIT_ROUTE } from './constants';

interface AutocompleteRequest {
  readonly prefix: string;
  readonly size: number;
}

// TODO: Add other fields
interface SimpleTrait {
  readonly id: number;
  readonly name: string;
}

class AutocompleteService {
  private requestor: Requestor

  constructor(requestor: Requestor) {
    this.requestor = requestor;
  }

  async autocompleteSimpleTrait(prefix: string, size: number): Promise<Array<SimpleTrait>> {
    const req: AutocompleteRequest = { prefix, size };
    const res = await this.requestor.post(AUTOCOMPLETE_SIMPLE_TRAIT_ROUTE, req);
    return res
  }
}

const autocompleteService = new AutocompleteService(requestor);
export default autocompleteService;
