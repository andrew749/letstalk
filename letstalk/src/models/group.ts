export interface Group {
  readonly groupId: number;
  readonly groupName: string;
}

export interface UserGroup extends Group {
  readonly id: number;
}
