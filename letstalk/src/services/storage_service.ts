import {FileSystem} from 'expo';
/**
 * A basic persistent key value store for the application.
 */
export interface StorageService {
  getKey(): Promise<any>;
}
