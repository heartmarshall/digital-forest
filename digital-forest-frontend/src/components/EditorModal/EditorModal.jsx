// src/components/EditorModal/EditorModal.jsx
import React from 'react';
import './EditorModal.css';
import AddPlantForm from '../AddPlantForm/AddPlantForm';

function EditorModal({ isOpen, onClose, onPlantAdded }) {
  if (!isOpen) {
    return null;
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close-button" onClick={onClose}>×</button>
        <h2>Нарисуй свое растение</h2>
        <AddPlantForm onPlantAdded={onPlantAdded} />
      </div>
    </div>
  );
}

export default EditorModal;