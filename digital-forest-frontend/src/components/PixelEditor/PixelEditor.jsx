// src/components/PixelEditor/PixelEditor.jsx
import React from 'react';
import './PixelEditor.css';

function PixelEditor({ 
  grid, 
  gridSize,
  colors,
  currentColor, 
  setCurrentColor,
  onCellClick 
}) {
  const [isMouseDown, setIsMouseDown] = React.useState(false);

  return (
    <div className="editor-container">
      <div className="palette">
        {colors.map((color, index) => (
          <button 
            key={index}
            className={`palette-color ${currentColor === color ? 'selected' : ''}`}
            style={{ backgroundColor: color === 'transparent' ? undefined : color }}
            {...(color === 'transparent' && { className: `palette-color eraser ${currentColor === color ? 'selected' : ''}` })}
            onClick={() => setCurrentColor(color)}
          />
        ))}
      </div>
      <div 
        className="grid" 
        style={{ gridTemplateColumns: `repeat(${gridSize}, 1fr)` }}
        onMouseDown={() => setIsMouseDown(true)}
        onMouseUp={() => setIsMouseDown(false)}
        onMouseLeave={() => setIsMouseDown(false)}
      >
        {grid.map((color, index) => (
          <div 
            key={index}
            draggable="false" // Запрещаем перетаскивание элемента
            className="grid-cell"
            style={{ backgroundColor: color }}
            onMouseDown={(event) => {
              event.preventDefault(); // Отменяем стандартное поведение браузера
              onCellClick(index);
            }}
            onMouseEnter={() => {
              if (isMouseDown) {
                onCellClick(index);
              }
            }}
          />
        ))}
      </div>
    </div>
  );
}

export default PixelEditor;