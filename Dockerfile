# Run the app in the environment
FROM letstalk_env as app
CMD ./run_remote.sh
EXPOSE 3000
