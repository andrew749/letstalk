import Immutable from 'immutable';

export interface Group {
  readonly groupId: number;
  readonly groupName: string;
}

export interface UserGroup extends Group {
  readonly id: number;
}

const groupNamesByGroupId: Immutable.Map<string, string> = Immutable.Map([
  ['WICS', 'Women in Computer Science'],
]);

export { groupNamesByGroupId };
