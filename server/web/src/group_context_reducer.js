const initialState = {
    currentGroup: undefined,
}

export const NEW_CURRENT_GROUP = "NEW_CURRENT_GROUP";

export function setNewCurrentGroupAction(group) {
    return {type: NEW_CURRENT_GROUP, group: group};
}

export function getCurrentGroup(state) {
    return state.groupContextReducer.group;
}

export function groupContextReducer(state = initialState, action) {
    switch(action.type) {
        case NEW_CURRENT_GROUP:
            return Object.assign({}, state, {group: action.group});
        default:
            return state;
    }
}

export class GroupContextManager {
    constructor() {
        this.newGroupActionListeners = [];
        this.checkGroupDidChange = this.checkGroupDidChange.bind(this);
        this.lastObservedGroup = undefined;
        this.subscribeListenerToNewGroupAction = this.subscribeListenerToNewGroupAction.bind(this);
    }

    /**
     * 
     * @param {*} listener A callback that will get called when the current group changes. Passed the state and new group 
     */
    subscribeListenerToNewGroupAction(listener) {
        this.newGroupActionListeners.push(listener);
    }

    getLastObserved() {
        return this.lastObservedGroup;
    }

    setLastObserved(newValue) {
        this.lastObservedGroup = newValue;
    }

    checkGroupDidChange(state) {
        const newGroup = getCurrentGroup(state);
        const lastObserved = this.getLastObserved();
        this.setLastObserved(newGroup)
        if (lastObserved && newGroup && lastObserved.groupId != newGroup.groupId || !lastObserved && !!newGroup) {
            console.log("Found change in group. Updating listeners");
            this.newGroupActionListeners.forEach(listener => listener(state, newGroup));
        }
    }
}
