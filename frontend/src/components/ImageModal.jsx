import { useState, useEffect } from 'react'
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Modal from '@mui/material/Modal';
import Typography from '@mui/material/Typography';

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 400,
    bgcolor: 'background.paper',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
  };

// TODO: fetch high res image when opened

function ImageModal({ open, setOpen, id, lowResImage, description, backendURL }) {
    
    const [image, setImage] = useState("")
    useEffect(() => {
        if (id !== "") {
            fetch(backendURL + "images/" + id + "?isHighResolution")
            .then(response => response.json())
            .then(jsonData => {
                setImage(jsonData.data)
            })
        }
      }, [id])
    
    const handleDownload = () => {
        // Convert base64 to blob
        const byteCharacters = atob(image);
        const byteNumbers = new Array(byteCharacters.length);
        for (let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i);
        }
        const byteArray = new Uint8Array(byteNumbers);
        const blob = new Blob([byteArray], {type: 'image/png'});

        // Create a download link and click it programmatically
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'image.png'); // or any other extension
        document.body.appendChild(link);
        link.click();

        // Clean up and revoke the URL
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
    };
    
    return (
        <>
        <Modal
            open={open}
            onClose={() => {
                setOpen(false)
                setImage("")
            }}
            aria-labelledby="modal-modal-title"
            aria-describedby="modal-modal-description"
        >
            <Box sx={style}>
                <img 
                    src={image === "" ? `data:image/jpeg;base64,${lowResImage}` : `data:image/jpeg;base64,${image}`}
                    width="100%"
                    />
                <Typography sx={{ fontSize: 14 }} color="text.secondary" gutterBottom>
                    {description}
                </Typography>
                <Button 
                    variant="contained"
                    size="large"
                    sx={{ margin: 3 }} 
                    onClick={() => handleDownload()}
                >
                    Download
                </Button>
            </Box>
        </Modal>
        </>
    );
}

export default ImageModal