import { useEffect, useState } from 'react'
import './App.css'
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Modal from '@mui/material/Modal';

import ImageGrid from './components/ImageGrid'
import ImageModal from './components/ImageModal'
import GenerateModal from './components/GenerateModal';

const backendURL = "http://localhost:8090/"

function App() {
    const [imageIds, setImageIds] = useState([])
    useEffect(() => {
        fetch(backendURL + "images")
        .then(response => response.json())
        .then(ids => setImageIds(ids))
    }, [])


    const [imageModalOpen, setImageModalOpen] = useState(false)
    const [focusedId, setFocusedId] = useState("")
    const [focusedImage, setFocusedImage] = useState("")
    const [focusedDescription, setFocusedDescription] = useState("")
    
    function setImageFocused(id, image, description) {
        setImageModalOpen(true)
        setFocusedId(id)
        setFocusedImage(image)
        setFocusedDescription(description)
    }


    const [generationModalOpen, setGenerationModalOpen] = useState(false)
    const [isLoading, setIsLoading] = useState(false)
    

    const [stylePrompt, setStylePrompt] = useState("")
    const [styleUrl, setStyleUrl] = useState("")
    const [scenePrompt, setScenePrompt] = useState("")
    const [sceneUrl, setSceneUrl] = useState("")
   
    function getConfig(prompt, url) {
        const value = prompt === "" ? url : prompt
        const isUrl = prompt === ""
        return [value, isUrl]
    }

    async function generateImage(numImages, size) {
        const [style, styleIsUrl] = getConfig(stylePrompt, styleUrl)
        const [scene, sceneIsUrl] = getConfig(scenePrompt, sceneUrl)
       

        setGenerationModalOpen(false)
        setIsLoading(true)
        const response = await fetch(backendURL + "images",
        {
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify({
                "style": style, 
                "styleIsUrl": styleIsUrl,
                "scene": scene,
                "sceneIsUrl": sceneIsUrl,
                "size": size,
                "numImages": numImages
            })
        })
        
        const newIds = await response.json();
        setImageIds(newIds.concat(imageIds))
        setIsLoading(false)
    }


    return (
        <>
        <Button 
            variant="contained"
            size="large"
            sx={{ margin: 3 }} 
            onClick={() => setGenerationModalOpen(true)}
        >
            Generate New Image
        </Button>
        <ImageGrid
            ids={imageIds}
            backendURL={backendURL}
            setImageFocused={setImageFocused}
        />
        <ImageModal
            open={imageModalOpen}
            setOpen={setImageModalOpen}
            id={focusedId}
            lowResImage={focusedImage}
            description={focusedDescription}
            backendURL={backendURL}
        />
        <GenerateModal
            open={generationModalOpen}
            setOpen={setGenerationModalOpen}
            generateImage={generateImage}
            setStylePrompt={setStylePrompt}
            setStyleUrl={setStyleUrl}
            setScenePrompt={setScenePrompt}
            setSceneUrl={setSceneUrl}
        />
        <Modal open={isLoading}>
                <Box sx={{ width: 50 }} className="modal-style">
                    <CircularProgress/>
                </Box>
            </Modal>
        </>
    )
}

export default App
