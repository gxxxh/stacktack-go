package worker

import (
	"fmt"
	"github.com/gxxxh/stacktack-go/src/worker"
	"github.com/tidwall/gjson"
	"testing"
	"time"
)

func TestNovaConsumerRun(t *testing.T) {
	config := worker.NewConfigByJson("E:\\gopath\\src\\stacktack-go\\test\\worker\\config.json")

	novaConsumer, err := worker.NewNovaConsumer(config)
	if err != nil {
		t.Error(err)
		return
	}
	novaConsumer.Run(worker.NovaNotificationInfoHandler)

}

func TestNovaConsumerShutDown(t *testing.T) {
	config := worker.NewConfigByJson("E:\\gopath\\src\\stacktack-go\\test\\worker\\config.json")

	novaConsumer, err := worker.NewNovaConsumer(config)
	if err != nil {
		t.Error(err)
		return
	}
	go novaConsumer.Run(worker.NovaNotificationInfoHandler)
	time.Sleep(10 * time.Second) //make sure shutdown should be called after Run
	novaConsumer.Shutdown(novaConsumer.ConsumerTag)
}

func TestHandleJson(t *testing.T) {
	data := "{\"_context_domain\": null, \"_context_request_id\": \"req-5db3a289-748d-4fbd-9feb-7dd99ceb85f4\", \"_context_global_request_id\": null, \"_context_quota_class\": null, \"event_type\": \"compute.instance.delete.start\", \"_context_service_catalog\": [{\"endpoints\": [{\"adminURL\": \"http://133.133.135.136:8778/placement\", \"region\": \"RegionOne\", \"internalURL\": \"http://133.133.135.136:8778/placement\", \"publicURL\": \"http://133.133.135.136:8778/placement\"}], \"type\": \"placement\", \"name\": \"placement\"}, {\"endpoints\": [{\"adminURL\": \"http://133.133.135.136:9292\", \"region\": \"RegionOne\", \"internalURL\": \"http://133.133.135.136:9292\", \"publicURL\": \"http://133.133.135.136:9292\"}], \"type\": \"image\", \"name\": \"glance\"}, {\"endpoints\": [{\"adminURL\": \"http://133.133.135.136:8776/v3/aac94320146c464ab84146e35aa61c77\", \"region\": \"RegionOne\", \"internalURL\": \"http://133.133.135.136:8776/v3/aac94320146c464ab84146e35aa61c77\", \"publicURL\": \"http://133.133.135.136:8776/v3/aac94320146c464ab84146e35aa61c77\"}], \"type\": \"volumev3\", \"name\": \"cinderv3\"}, {\"endpoints\": [{\"adminURL\": \"http://133.133.135.136:9696\", \"region\": \"RegionOne\", \"internalURL\": \"http://133.133.135.136:9696\", \"publicURL\": \"http://133.133.135.136:9696\"}], \"type\": \"network\", \"name\": \"neutron\"}], \"timestamp\": \"2023-03-17 06:59:36.321294\", \"_context_user\": \"f8db2401acfb4c3b98400dac8fa22207\", \"_unique_id\": \"974aa0d32a2e4c73a1b716dc79fd9215\", \"_context_resource_uuid\": null, \"_context_is_admin_project\": true, \"_context_read_deleted\": \"no\", \"_context_user_id\": \"f8db2401acfb4c3b98400dac8fa22207\", \"payload\": {\"state_description\": \"deleting\", \"availability_zone\": null, \"terminated_at\": \"\", \"ephemeral_gb\": 0, \"instance_type_id\": 1, \"deleted_at\": \"\", \"reservation_id\": \"r-dtkzx0j0\", \"instance_id\": \"b5762ff2-2e15-4aac-a26a-0381f618f1b5\", \"display_name\": \"1\", \"hostname\": \"1\", \"state\": \"error\", \"progress\": \"\", \"launched_at\": \"\", \"metadata\": {}, \"node\": null, \"ramdisk_id\": \"\", \"access_ip_v6\": null, \"disk_gb\": 1, \"access_ip_v4\": null, \"kernel_id\": \"\", \"host\": null, \"user_id\": \"f8db2401acfb4c3b98400dac8fa22207\", \"image_ref_url\": \"http://133.133.135.136:9292/images/\", \"cell_name\": \"\", \"root_gb\": 1, \"tenant_id\": \"aac94320146c464ab84146e35aa61c77\", \"created_at\": \"2023-03-17 06:46:24+00:00\", \"memory_mb\": 512, \"instance_type\": \"m1.tiny\", \"vcpus\": 1, \"image_meta\": {\"min_disk\": \"1\", \"container_format\": \"bare\", \"min_ram\": \"0\", \"disk_format\": \"qcow2\", \"base_image_ref\": \"\"}, \"architecture\": null, \"os_type\": null, \"instance_flavor_id\": \"1\"}, \"_context_project_name\": \"admin\", \"_context_system_scope\": null, \"_context_user_identity\": \"f8db2401acfb4c3b98400dac8fa22207 aac94320146c464ab84146e35aa61c77 - default default\", \"_context_auth_token\": \"gAAAAABkFA_V5wB71aB-17QYitx1V7tMzvJ5Bf-Nos52wXJs_Zcts5QtjYJK3ie-kFxbJSlO2I3JB1EiQwXguDug9akqFZ1_pgjlPiNzPl0BrhPt5uux6qCwB_zOAeDzqDgW2reBZ9swT_B9lgk2GPy7Ln4_FTDI1Z8HNG_MZ9CxptjBf-BD25WJW3QuLU345jVxjaDEQWdu\", \"_context_show_deleted\": false, \"_context_tenant\": \"aac94320146c464ab84146e35aa61c77\", \"_context_roles\": [\"reader\", \"admin\", \"member\"], \"priority\": \"INFO\", \"_context_read_only\": false, \"_context_is_admin\": true, \"_context_project_id\": \"aac94320146c464ab84146e35aa61c77\", \"_context_project_domain\": \"default\", \"_context_timestamp\": \"2023-03-17T06:59:36.129223\", \"_context_user_domain\": \"default\", \"_context_user_name\": \"admin\", \"publisher_id\": \"compute.openstack\", \"message_id\": \"2ced7ca3-3b46-4a90-abe1-a0ec772e8ad2\", \"_context_project\": \"aac94320146c464ab84146e35aa61c77\", \"_context_remote_address\": \"133.133.135.136\"}"
	jsonBody := gjson.Parse(data)
	fmt.Println(jsonBody.Get("event_type").String())
	payLoad := jsonBody.Get("payload")
	instanceID := ""
	if payLoad.Get("instance_id").Exists() {
		instanceID = payLoad.Get("instance_id").String()
	} else if payLoad.Get("instance_uuid").Exists() {
		instanceID = payLoad.Get("instacnce_uuid").String()
	} else if payLoad.Get("exception.kwargs.uuid").Exists() {
		instanceID = payLoad.Get("exception.kwargs.uuid").String()
	} else if payLoad.Get("instance.uuid").Exists() {
		instanceID = payLoad.Get("instance.uuid").String()
	}
	fmt.Println(instanceID)
}
