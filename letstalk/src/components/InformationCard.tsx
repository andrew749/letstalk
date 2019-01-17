import React, { ReactElement, ReactNode } from 'react';
import { AsyncStorage, View, Text, StyleSheet, Image, TouchableOpacity } from 'react-native';
import Colors from '../services/colors';
import Card from '../components/Card';
import { Button } from '../components';
import { Linking } from 'expo';
import navService from '../services/navigation-service';
import Collapsible from 'react-native-collapsible';
import { MaterialIcons } from '@expo/vector-icons';


enum InformationCardType {
    CLUB_DAY = 'information-card-club-day-visibilitya',
    PROFILE_FILL_CALL_TO_ACTION = 'profile-fill-ctaaaa',
}

interface Props {
    title: string;
    cardType: InformationCardType;
    onDismiss?: () => void;
}

interface State {
    informationCardVisibility: InformationCardVisibilityState;
}

enum InformationCardVisibilityState {
    VISIBLE = "visible",
    INVISIBLE = "invisible",
}

/**
 * Generic information card which we can get from the server.
 */
class InformationCard extends React.Component<Props, State> {
    state = {
        informationCardVisibility: InformationCardVisibilityState.VISIBLE
    }

    constructor(props: Props) {
        super(props);
        this.onCancel = this.onCancel.bind(this);
    }

    async componentDidMount() {
        const visibility = await AsyncStorage.getItem(this.props.cardType);
        this.setState({ informationCardVisibility: visibility == null ? InformationCardVisibilityState.VISIBLE : visibility as InformationCardVisibilityState })
    }

    public async onCancel() {
        dismissCard(this.props.cardType, () => this.setState({ informationCardVisibility: InformationCardVisibilityState.INVISIBLE }));
        this.props.onDismiss();
    }

    shouldRenderComponent = () => {
        return this.state.informationCardVisibility == InformationCardVisibilityState.VISIBLE;
    }

    render() {
        const shouldRender = this.shouldRenderComponent();
        return shouldRender ?
            <Card style={styles.cardOverrides}>
                <TouchableOpacity style={styles.cancelButtonContainer} onPress={this.onCancel}>
                    <MaterialIcons name="close" size={18} />
                </TouchableOpacity>
                <Text style={[styles.cardHeader]}>
                    {this.props.title}
                </Text>
                <View>
                    {this.props.children}
                </View>
            </Card>
            : <View></View>;
    }
}

async function dismissCard(cardType: InformationCardType, action: () => void) {
    await AsyncStorage.setItem(cardType, InformationCardVisibilityState.INVISIBLE);
}

interface ClubDayProps { }
interface ClubDayState {
    isCollapsed: boolean;
}

export class ClubDayInformationCard extends React.Component<ClubDayProps, ClubDayState> {
    constructor(props: Props) {
        super(props);
        this.state = {
            isCollapsed: true,
        };
    }

    render() {
        const moreButton = (this.state.isCollapsed) ? <TouchableOpacity onPress={() => this.setState({ isCollapsed: false })}>
            <Text style={styles.moreButton} >Show more</Text>
        </TouchableOpacity> : undefined;
        return (
            <InformationCard title="Welcome to the Hive!" cardType={InformationCardType.CLUB_DAY}>
                <Text style={[styles.textSection]}>
                    Matches will be coming out in the next couple of weeks. Stay tuned!
                </Text>
                <Text style={styles.textSection}>
                    In the meanwhile, feel free to search for people you might be interested in connecting with. For example:
                </Text>
                {moreButton}
                <Collapsible collapsed={this.state.isCollapsed}>
                    <Text style={[styles.points]}>- Meet other people in your cohort</Text>
                    <Text style={[styles.points]}>- Ask for tips from a person who worked at your dream company</Text>
                    <Text style={[styles.points]}>- Find someone who went on an exchange term</Text>
                    <Text style={[styles.textSection, { marginTop: 20 }]}>Expand your horizons. Grow your network.</Text>
                    <Text style={[styles.signature]}>The Hive Team</Text>
                    <View style={styles.imageContainer}>
                        <Image style={styles.imageStyle} source={require('../img/logo_android.png')} />
                    </View>
                </Collapsible>
            </InformationCard>
        );
    }
}

interface ProfileFillCallToActionProps { }
interface ProfileFillCallToActionState { 
    isHidden: boolean;
}

export class ProfileFillCallToAction extends React.Component<ProfileFillCallToActionProps, ProfileFillCallToActionState> {
    constructor(props: Props) {
        super(props);
        this.state = {isHidden: false};
    }

    render() {
        return (this.state.isHidden) ? <View/> : (
            <InformationCard title="Tell us about yourself" cardType={InformationCardType.PROFILE_FILL_CALL_TO_ACTION} onDismiss={() => this.setState({ isHidden: true })}>
                <Text style={[styles.textSection]}>
                    Your profile section gives others a first glance into your world. The more you give, the easier it'll be for others to relate to your experiences and understand you. 
                </Text>
                <Text style={[styles.textSection]}>
                    Share your favorite spot in Waterloo, a hobby you're working on, or something you know way too much about!
                </Text>
                <Button
                    buttonStyle={{ backgroundColor: Colors.HIVE_PRIMARY }}
                    textStyle={{ color: Colors.WHITE }}
                    onPress={async () => {
                        await dismissCard(InformationCardType.PROFILE_FILL_CALL_TO_ACTION, () => { });
                        this.setState({isHidden: true});
                        await navService.navigate("Profile", {});
                    }}
                    title="Go to your profile"
                />
            </InformationCard>
        );
    }
}

const PADDING_TEXT_SECTION = 10;

const styles = StyleSheet.create({
    moreButton: {
        color: '#BDBDBD',
    },
    cancelButtonContainer: {
        justifyContent: 'flex-end',
        alignItems: 'flex-end',
    },
    cardOverrides: {
        marginHorizontal: 0,
    },
    textSection: {
        fontSize: 18,
        paddingTop: PADDING_TEXT_SECTION,
        paddingBottom: PADDING_TEXT_SECTION,
    },
    cardHeader: {
        fontSize: 24,
        fontWeight: 'bold',
    },
    imageStyle: {
        marginTop: 20,
        width: 100,
        height: 100,
        alignItems: 'center',
    },
    imageContainer: {
        justifyContent: 'center',
        alignItems: 'center',
    },
    signature: {
        fontSize: 18,
        marginTop: 20,
        fontWeight: 'bold',
    },
    points: {
        fontSize: 16
    },
});