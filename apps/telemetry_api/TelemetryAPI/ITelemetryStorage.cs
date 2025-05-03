using TelemetryAPI.Consumers;

public interface ITelemetryStorage
{
    public void Add(SensorUpdated sensor);
    public IEnumerable<SensorUpdated> GetAll();
    public IEnumerable<SensorUpdated> GetBySensorId(string sensorId);
}