# make startProm
# this file if you want run this without running a docker-compose.yml
.PHONY: startProm
startProm:
	docker run \
	--rm \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

# make startGrafana
# for first timers, the username & password is both `admin`
.PHONY: startGrafana
startGrafana:
	docker volume create grafana_data
	docker volume inspect grafana_data
	docker run -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana