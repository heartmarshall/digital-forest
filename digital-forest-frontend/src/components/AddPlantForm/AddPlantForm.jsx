// src/components/AddPlantForm/AddPlantForm.jsx
import React, { useState, useEffect } from 'react';
import PixelEditor from '../PixelEditor/PixelEditor';
import { gridToBase64 } from '../../utils/imageConverter';
import { createPlant } from '../../api/plantApi';
import './AddPlantForm.css';

const GRID_SIZE = 16;
const COLORS = ['transparent', '#FFFFFF', '#000000', '#FF0000', '#00FF00', '#0000FF', '#FFFF00', '#FF00FF', '#00FFFF'];
const initialGrid = Array(GRID_SIZE * GRID_SIZE).fill(COLORS[0]);

function AddPlantForm({ onPlantAdded }) {
  const [author, setAuthor] = useState('');
  const [grid, setGrid] = useState(initialGrid);
  const [currentColor, setCurrentColor] = useState(COLORS[2]);
  
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  // Сообщение об успехе нам больше не нужно, т.к. окно просто закроется
  
  const handleCellClick = (index) => {
    const newGrid = [...grid];
    newGrid[index] = currentColor;
    setGrid(newGrid);
  };

  const handleClearGrid = () => {
    setGrid(initialGrid);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!author.trim()) {
      setError('Пожалуйста, введите ваше имя или никнейм.');
      return;
    }
    setIsSubmitting(true);
    setError('');

    try {
      const imageData = await gridToBase64(grid, GRID_SIZE);
      await createPlant(author, imageData);
      
      // Сбрасываем форму
      setAuthor('');
      setGrid(initialGrid);

      // Вызываем колбэк, который закроет окно и обновит лес
      if (onPlantAdded) {
        onPlantAdded();
      }
    } catch (err) {
      setError('Не удалось добавить растение. Попробуйте снова.');
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    // Убираем внешний div, так как компонент теперь внутри модального окна
    <>
      <h2>Нарисуй свое растение</h2>
      <PixelEditor 
        grid={grid}
        gridSize={GRID_SIZE}
        colors={COLORS}
        currentColor={currentColor}
        setCurrentColor={setCurrentColor}
        onCellClick={handleCellClick}
      />
      <div className="editor-controls">
        <button type="button" onClick={handleClearGrid}>Очистить</button>
      </div>
      <form className="plant-form" onSubmit={handleSubmit}>
        <input
          type="text"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          placeholder="Ваше имя или никнейм"
          maxLength="255"
          required
          disabled={isSubmitting}
        />
        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Добавляем...' : 'Добавить в лес'}
        </button>
      </form>
      {error && <p className="error-message">{error}</p>}
    </>
  );
}

export default AddPlantForm;