use spiffe::workload_api::client::WorkloadApiClient;
use spiffe::spiffe_id::{SpiffeId};

fn main() {
    let client = WorkloadApiClient::new("unix:/run/iotedge/sockets/workloadapi.sock").unwrap();
    let spiffe_id = SpiffeId::try_from("spiffe://iotedge/my-service").unwrap();

    let target_audience = &["spiffe://iotedge/mqttbroker"];
    // fetch a jwt token for the provided SPIFFE-ID and with the target audience `service1.com`
    let jwt_svid = client.fetch_jwt_svid(target_audience, Some(&spiffe_id)).unwrap();
    let test = jwt_svid.spiffe_id().path().clone();
    println!("{}", test.to_string());
}
