import { NavigationScreenProp, NavigationStackAction } from "react-navigation";
import { Font } from "expo";
import { AsyncStorage } from "react-native";

class TutorialService {
    static SHOW_TUTORIAL_KEY = "tutorial";

    static async shouldShowTutorial(): Promise<boolean> {
      const tutorialState = await AsyncStorage.getItem(this.SHOW_TUTORIAL_KEY);
      return tutorialState !== "true";
    }

    static async ackShowTutorial(): Promise<void> {
      await AsyncStorage.setItem(this.SHOW_TUTORIAL_KEY, "true")
    }

    static async launchTutorial(navigation: NavigationScreenProp<void, NavigationStackAction>) {
      const shouldShowTutorial = await this.shouldShowTutorial();
      if (!shouldShowTutorial) return;

      // load for the walkthrough
      await Font.loadAsync({
        'Arial': require('../assets/fonts/Arial.ttf'),
      });

      await navigation.push('WalkthroughView');
    }
}

export default TutorialService;
