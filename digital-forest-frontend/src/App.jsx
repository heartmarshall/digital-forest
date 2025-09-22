// src/App.jsx
import { useState, useEffect } from 'react';
import Forest from './components/Forest/Forest';
import EditorModal from './components/EditorModal/EditorModal';
import { fetchRandomPlants } from './api/plantApi';
import './App.css';

function App() {
  const [plants, setPlants] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isEditorOpen, setIsEditorOpen] = useState(false);

  const loadPlants = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchRandomPlants(50);
      setPlants(data);
    } catch (err) {
      setError("Не удалось загрузить цифровой лес. Попробуйте позже.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlants();
  }, []);

  const handlePlantAdded = () => {
    setIsEditorOpen(false);
    loadPlants();
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <div className="header-title">
          <h1>Цифровой лес</h1>
          <p>Описание проекта</p>
        </div>
        <button className="add-plant-button" onClick={() => setIsEditorOpen(true)}>
          Добавить растение
        </button>
      </header>
      
      <main>
        <Forest plants={plants} loading={loading} error={error} />
      </main>
      
      <EditorModal 
        isOpen={isEditorOpen}
        onClose={() => setIsEditorOpen(false)}
        onPlantAdded={handlePlantAdded}
      />
    </div>
  )
}

export default App;