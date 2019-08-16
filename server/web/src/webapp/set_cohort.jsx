import React from 'react';
import '../scss/group_register_page.scss';
import {Alert, Container, Form, Button} from 'react-bootstrap';
import ClipLoader from 'react-spinners/ClipLoader';

import apiServiceConnect from '../api/api_service_connect';
import {getFetchInfo} from '../api/api_module';
import {getCohortsApiModule} from '../api/get_cohorts_module';

const MENTORSHIP_PREFERENCE_MENTOR = 1;
const MENTORSHIP_PREFERENCE_MENTEE = 2;
const MENTORSHIP_PREFERENCE_NONE = 3;


class SetCohortPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            programId: "",
            sequenceId: "",
            gradYear: -1,
            mentorshipPreference: -1,
            bio: "",
            hometown: "",
        };

        this.handleChange = this.handleChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.handleChangeInt = this.handleChangeInt.bind(this);
    }

    validateForm() {
        return this.state.programId.length > 0
            && this.state.sequenceId.length > 0
            && this.state.gradYear.length !== -1  // Change in 9999
            && this.state.mentorshipPreference !== -1
            ;
    }

    componentDidMount() {
        this.props.getCohorts();
    }

    handleChange(event) {
        this.setState({
            [event.target.id]: event.target.value
        });
    }

    handleChangeInt(event) {
        this.setState({
            [event.target.id]: parseInt(event.target.value, 10)
        });
    }

    onSubmit(event) {
        event.preventDefault();
        console.log(this.state)
    }

    renderProgramOptions() {
        const { cohorts } = this.props;
        if (!cohorts) return [];
        let programIds = cohorts.map(cohort => cohort.programId)
        programIds = [...new Set(programIds)];
        return programIds.map(programId => {
            return <option key={programId} value={programId}>{programId}</option>;
        });
    }

    renderSequenceOptions() {
        const { cohorts } = this.props;
        if (!cohorts) return [];
        const { programId } = this.state;
        let sequenceIds = cohorts.filter(cohort => {
          return programId.length === 0 || cohort.programId === programId
        }).map(cohort => cohort.sequenceId)
        sequenceIds = [...new Set(sequenceIds)];
        return sequenceIds.map(sequenceId => {
            return <option key={sequenceId} value={sequenceId}>{sequenceId}</option>;
        });
    }

    renderGradYearOptions() {
        const { cohorts } = this.props;
        if (!cohorts) return [];
        const { programId, sequenceId } = this.state;
        let gradYears = cohorts.filter(cohort => {
          return (programId.length === 0 || cohort.programId === programId) &&
            (sequenceId.length === 0 || cohort.sequenceId === sequenceId);
        }).map(cohort => cohort.gradYear)
        gradYears = [...new Set(gradYears)];
        return gradYears.map(gradYear => {
            return <option key={gradYear} value={gradYear}>{gradYear}</option>;
        });
    }

    render() {
        const { fetchState, errorMessage } = this.props.getCohortsFetchInfo;
        let body = null;
        if (fetchState === "error") {
            body = (
                <Alert variant="danger">
                    Failed to load cohorts with error: {errorMessage}
                </Alert>
            );
        } else if (fetchState === "fetching" || fetchState === 'prefetch') {
            body = <ClipLoader />;
        } else {
            body = (
                <Form onSubmit={this.onSubmit}>
                    <Form.Group controlId="programId">
                        <Form.Label>Program</Form.Label>
                        <Form.Control
                            onChange={this.handleChange}
                            as="select">
                            <option key={"-1"} value={""}>Choose one</option>
                            {this.renderProgramOptions()}
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="sequenceId">
                        <Form.Label>Co-op Sequence</Form.Label>
                        <Form.Control
                            onChange={this.handleChange}
                            as="select">
                            <option value={""}>Choose one</option>
                            {this.renderSequenceOptions()}
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="gradYear">
                        <Form.Label>Graduating Year</Form.Label>
                        <Form.Control
                            onChange={this.handleChangeInt}
                            as="select">
                            <option value={-1}>Choose one</option>
                            {this.renderGradYearOptions()}
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="mentorshipPreference">
                        <Form.Label>Your Preferred Role</Form.Label>
                        <Form.Control
                            onChange={this.handleChangeInt}
                            as="select">
                            <option value={-1}>Choose one</option>
                            <option value={MENTORSHIP_PREFERENCE_MENTOR}>Mentor</option>
                            <option value={MENTORSHIP_PREFERENCE_MENTEE}>Mentee</option>
                            <option value={MENTORSHIP_PREFERENCE_NONE}>Okay with either</option>
                            <option value={MENTORSHIP_PREFERENCE_NONE}>I don't know</option>
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="bio">
                        <Form.Label>Bio</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.bio}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="hometown">
                        <Form.Label>Hometown</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.hometown}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Button
                        block
                        disabled={!this.validateForm()}
                        type="submit"
                    >
                        Submit
                    </Button>
                </Form>
            );
        }

        return (
            <Container className="panel-body">
                {body}
            </Container>
        );
    }
}

const SetCohortPageComponent = apiServiceConnect(
    (state) => {
        return {
            cohorts: getCohortsApiModule.getData(state),
            getCohortsFetchInfo: getFetchInfo(getCohortsApiModule, state),
        }
    },
    (dispatch) => {
        return {
            getCohorts: () => {
                return dispatch(getCohortsApiModule.getApiExecuteAction({}))
            },
        }
    }
)(SetCohortPage);

export default SetCohortPageComponent;
