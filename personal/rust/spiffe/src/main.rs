use spiffe::workload_api::client::WorkloadApiClient;

fn main() {
    let client = WorkloadApiClient::new("unix:/run/iotedge/sockets/workloadapi.sock").unwrap();

    let jwt_bundles_set = client.fetch_jwt_bundles().unwrap();

    println!("{:?}", jwt_bundles_set);

    let target_audience = &["spiffe://iotedge/mqttbroker"];
    // fetch a jwt token for the provided SPIFFE-ID and with the target audience `service1.com`
    let jwt_svid = client.fetch_jwt_svid(target_audience, None).unwrap();
    let test = jwt_svid.spiffe_id().clone();
    println!("{:?}", test);
}
