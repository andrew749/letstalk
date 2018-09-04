import { NavigationContainerComponent, NavigationActions } from "react-navigation";
/**
 * THIS IS COPIED VERBATIM FROM
 * https://reactnavigation.org/docs/en/navigating-without-navigation-prop.html
 */


class NavigationService {
  private _navigator: NavigationContainerComponent;
  setTopLevelNavigator(navContainer: NavigationContainerComponent) {
    this._navigator = navContainer;
  }

  navigate(routeName: string, params: any){
    console.log("navigating to " + routeName)
    this._navigator.dispatch(
      NavigationActions.navigate({
        routeName,
        params,
      })
    );
  }
}

let navService = new NavigationService();

export default navService;
