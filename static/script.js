const MAX_POINTS = 25;
const tempData = [];
const humData = [];
const labels = [];

const ctx = document.getElementById('sensorChart').getContext('2d');
const chart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: labels,
    datasets: [
      {
        label: 'Temperatura (°C)',
        data: tempData,
        borderColor: '#ef4444',
        backgroundColor: 'rgba(239, 68, 68, 0.08)',
        borderWidth: 2.5,
        pointRadius: 2,
        pointHoverRadius: 5,
        tension: 0.3,
        fill: true,
        cubicInterpolationMode: 'monotone'
      },
      {
        label: 'Umidade (%)',
        data: humData,
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59, 130, 246, 0.08)',
        borderWidth: 2.5,
        pointRadius: 2,
        pointHoverRadius: 5,
        tension: 0.3,
        fill: true,
        cubicInterpolationMode: 'monotone'
      }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    animation: {
      duration: 300,
      easing: 'easeOutQuart'
    },
    scales: {
      x: {
        grid: { display: false },
        ticks: {
          maxRotation: 0,
          autoSkip: true,
          maxTicksLimit: 8,
          color: '#64748b'
        }
      },
      y: {
        beginAtZero: false,
        min: 0,
        max: 100,
        grid: {
          color: 'rgba(0,0,0,0.04)'
        },
        ticks: {
          color: '#64748b',
          stepSize: 10
        }
      }
    },
    plugins: {
      legend: {
        labels: {
          color: '#1e293b',
          font: { size: 13, weight: '500' },
          usePointStyle: true,
          padding: 20
        }
      },
      tooltip: {
        mode: 'index',
        intersect: false,
        backgroundColor: 'rgba(30, 41, 59, 0.9)',
        titleFont: { size: 13 },
        bodyFont: { size: 13 },
        padding: 12,
        callbacks: {
          label: function(context) {
            let label = context.dataset.label || '';
            if (label) label += ': ';
            label += context.parsed.y.toFixed(1);
            if (context.dataset.label.includes('Temperatura')) label += ' °C';
            else if (context.dataset.label.includes('Umidade')) label += ' %';
            return label;
          }
        }
      }
    },
    interaction: {
      mode: 'nearest',
      axis: 'x',
      intersect: false
    }
  }
});

const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  const now = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  console.log(data);
  document.getElementById('temp-value').textContent = data.temperatura.toFixed(1);
  document.getElementById('hum-value').textContent = data.umidade.toFixed(1);
  document.getElementById('sensor-id').textContent = data.ID;

  const controleEl = document.getElementById('controle-indicator');
  const labelEl = document.getElementById('controle-label');
  
  
  if (data.controle === true) {
    controleEl.className = 'status-indicator active';
    labelEl.textContent = 'Ativo';
  } else if (data.controle === false) {
    controleEl.className = 'status-indicator inactive';
    labelEl.textContent = 'Inativo';
  } else {
    controleEl.className = 'status-indicator';
    labelEl.textContent = '--';
  }

  labels.push(now);
  tempData.push(data.temperatura);
  humData.push(data.umidade);

  if (labels.length > MAX_POINTS) {
    labels.shift();
    tempData.shift();
    humData.shift();
  }

  chart.update('none');
};

ws.onerror = () => {
  console.error('Falha na conexão WebSocket.');
};