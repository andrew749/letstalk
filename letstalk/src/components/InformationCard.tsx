import React from 'react';
import { View, Text, StyleSheet, Image } from 'react-native';
import Colors from '../services/colors';
import Card from '../components/Card';

interface Props {}

/**
 * Generic information card which we can get from the server.
 */
const InformationCard: React.SFC<Props> = props => {
    return (
        <Card style={styles.cardOverrides}>
            <Text style={[styles.cardHeader]}>
                Welcome to the Hive!!
            </Text>
            <Text style={[styles.textSection]}>
                Matches will be coming out in the next couple of weeks. Stay tuned!
            </Text>
            <Text style={styles.textSection}>
                In the meanwhile, feel free to search for people you might be interested in connecting with. For example:
            </Text>
            <Text style={[styles.points]}>- Meet other people in your cohort</Text>
            <Text style={[styles.points]}>- Ask for tips from a person who worked at your dream company</Text>
            <Text style={[styles.points]}>- Find someone who went on an exchange term</Text>
            <Text style={[styles.textSection, {marginTop: 20}]}>Expand your horizons. Grow your network.</Text>
            <Text style={[styles.signature]}>The Hive Team</Text>
            <View style={styles.imageContainer}>
                <Image style={styles.imageStyle} source={require('../img/logo_android.png')}/>
            </View>
        </Card>
    );
}

export default InformationCard;

const PADDING_TEXT_SECTION = 10;

const styles = StyleSheet.create({
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