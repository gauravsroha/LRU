import React, { useState } from 'react';
import './App.css';

function App() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [expiration, setExpiration] = useState(5);
  const [result, setResult] = useState('');

  const handleGet = async () => {
    try {
      const response = await fetch(`http://localhost:8080/get?key=${key}`);
      if (response.ok) {
        const data = await response.json();
        setResult(`Value: ${data.value}`);
      } else {
        setResult('Key not found');
      }
    } catch (error) {
      setResult('Error fetching data');
    }
  };

  const handleSet = async () => {
    try {
      const response = await fetch('http://localhost:8080/set', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key, value, expiration: parseInt(expiration) }),
      });
      if (response.ok) {
        const data = await response.json();
        setResult('Value set successfully: ' + data.status);
      } else {
        const errorText = await response.text();
        setResult('Error setting value: ' + errorText);
      }
    } catch (error) {
      setResult('Error setting data: ' + error.message);
    }
  };

  return (
    <div className="App">
      <h1>LRU Cache Assignment</h1>
      <div>
        <input
          type="text"
          placeholder="Key"
          value={key}
          onChange={(e) => setKey(e.target.value)}
        />
        <input
          type="text"
          placeholder="Value"
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />
        <input
          type="number"
          placeholder="Expiration (seconds)"
          value={expiration}
          onChange={(e) => setExpiration(e.target.value)}
        />
        <button onClick={handleSet}>Set</button>
        <button onClick={handleGet}>Get</button>
      </div>
      <div>
        <p>{result}</p>
      </div>
    </div>
  );
}

export default App;