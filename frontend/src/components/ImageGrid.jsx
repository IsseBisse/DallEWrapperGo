import Grid from '@mui/material/Grid';

import ImageThumbnail from './ImageThumbnail';

function ImageGrid({ ids, backendURL, setImageFocused }) {
    return (
        <Grid 
            container 
            spacing={{ xs: 2, md: 3 }} 
            columns={{ xs: 4, sm: 8, md: 12 }}
        >
            {ids.map(id => (
                <Grid item xs={2} sm={4} md={4} key={id}>
                    <ImageThumbnail 
                        id={id}
                        backendURL={backendURL}
                        setImageFocused={setImageFocused}
                    />
                </Grid>
            ))}
        </Grid>
    )
}

export default ImageGrid