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