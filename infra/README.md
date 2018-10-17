# Hive Infra Runbook
This is a runbook for various infra actions you may want to perform as oncall.

# Foreword/General Advice
I would recommend writing everything you do down, just in case you don't remember what state you put the server into. Could help with debugging in the future.

NEVER delete data in production unless another engineer is looking over your shoulder. We do take backups for a doomsday scenario, but restoring from a backup gracefully could prove to be very difficult.

# Core app structure
The main application is located at `/var/app/letstalk` on the host. Infrastructure related code and scripts to help with admin tasks are in `/var/app/letstalk/infra/`. Frontend related code is located in `/var/app/letstalk/landing`.

Currently, in production we have our backend server hosting an instance of nginx to perform ssl termination on the host. Both the application and nginx instance are run inside docker containers and connected via a virtual mesh network. This application server makes a db connection to a mysql RDS database which is hosted in another availability zone. Similarly, we host an elastic search instance in another availability zone in AWS that the backend server connects to.

Notifications are a little more complicated. Notifications are fed to SQS, under the queue Notifications, and then read by a lambda instance. The lambda instance has some security group magic with NATing (specifically Elastic NAT) and such to allow it to talk to the outside world.

# Monitoring
The following are some interesting monitoring metrics to see application health:

## Datadog
https://app.datadoghq.com/dash/902632/main-service?live=true&page=0&is_auto=false&from_ts=1539238396819&to_ts=1539324796819&tile_size=m
This dashboard includes the traffic that we are receiving and what codes we are returning. It also has a graph tracking if the instance is up. Useful for seeing the flow of good traffic vs errors.

## Sentry
https://sentry.io/hive-an/hive/
This shows the application errors that we are encountering (including both mobile and server side). It's important to note that errors are tagged with whether they occur in production or development, so filtering by an environment can help reduce false positives.

# Helpful Commands

## View docker logs
You need to be in `/var/app/letstalk` when running these commands.

### View all logs for all dockerized apps
```
docker-compose logs
```

### Tail all logs for all dockerized apps
```
docker-compose logs -f
```

### View logs for a specific dockerized app
```
docker-compose logs <app_name>
```

Where app_name can be determined from viewing which docker services are running.

For example:
```
docker-compose logs letstalk_webapp_1
```

## View which docker services are running
```
docker ps
```

## Connect to the database
```
cd /var/app/letstalk; ./connect_db.sh
```

## Deploy new app version

### Get new source code
```
git pull
```

### Deploy backend and frontend
Make sure you are in the infra folder!!!

```
cd /var/app/letstalk/infra;
```

Run deploy script:

```
./deploy.sh
```

## Deploy new frontend application
Once you have updated code in /var/app/letstalk/landing, you can run:

```
cd /var/app/letstalk/infra; ./build_frontend.sh; ./prod.sh
```

## Deploy new backend application

```
cd /var/app/letstalk/infra; ./build_prod.sh; ./prod.sh
```

## Clean up space (if the server is low on space)
This command removes stale images and containers from docker so that we can reclaim some space. If we are really tight on space, we can remove some of oldest db backups from /var (need sudo).
```
docker system prune
```

## Reload load balancer configuration
Sometimes you might need to reload nginx server. You can do this by:
```
cd /var/app/letstalk/infra/lb; ./reload_lb.sh
```

## Tear down server (BE CAREFUL)
The only real reason you might want this command is when renewing certs since this requires a small amount of downtime so that nginx can bind to the port to talk to letsencrypt.
To tear down the load balancer and container gracefully run:

```
cd /var/app/letstalk; ./bring_down.sh
```

