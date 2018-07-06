import { NavigationScreenProp, NavigationStackAction } from "react-navigation";
import { Font } from "expo";
import { AsyncStorage } from "react-native";

class AppStateService {
    static async isFirstLaunch(): Promise<boolean> {
      const tutorialState = await AsyncStorage.getItem("first_launch");
      // return tutorialState === "true";
      return true;
    }

    static async launchTutorialIfFirstLaunch(navigation: NavigationScreenProp<void, NavigationStackAction>) {
      const firstLaunch = await this.isFirstLaunch();
      if (firstLaunch) {
          await this.launchTutorial(navigation);
      }
    }

    static async launchTutorial(navigation: NavigationScreenProp<void, NavigationStackAction>) {
      // load for the walkthrough
      await Font.loadAsync({
        'Arial': require('../assets/fonts/Arial.ttf'),
      });

      navigation.push('WalkthroughView');
    }
}

export default AppStateService;
