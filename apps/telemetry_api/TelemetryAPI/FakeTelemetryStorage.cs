using TelemetryAPI.Consumers;

public class FakeTelemetryStorage : ITelemetryStorage
{
    public static List<SensorUpdated> sensorUpdateds = _sensorUpdateds ??= [];
    private static List<SensorUpdated>? _sensorUpdateds;

    public void Add(SensorUpdated sensor)
    {
        if (sensorUpdateds.Count > 100)
            sensorUpdateds.Remove(sensorUpdateds.First());

        sensorUpdateds.Add(sensor);
    }

    public IEnumerable<SensorUpdated> GetAll() => sensorUpdateds;
    public IEnumerable<SensorUpdated> GetBySensorId(string sensorId) => sensorUpdateds.Where(x => x.SensorId == sensorId);
}
