package model

type ProcessSchema struct {
	Model
	DoesTransform bool   `db:"does_transform" json:"does_transform"`
	Note          string `db:"note" json:"note"`
	ProcessType   string `db:"process_type" json:"process_type"`
	TemplateID    string `db:"template_id" json:"template_id"`
	TemplateName  string `db:"template_name" json:"template_name"`
}

type ProcessExtendedModel struct {
	ProcessSchema
	Setup         ProcessSetup         `json:"setup"`
	InputSamples  []ProcessSample      `json:"input_samples"`
	OutputSamples []ProcessSample      `json:"output_samples"`
	Measurements  []ProcessMeasurement `json:"measurements"`
}

type ProcessSetup struct {
}

type ProcessSample struct {
}

type ProcessMeasurement struct {
}

/*
       setup: r.table('process2setup').getAll(process('id'), {index: 'process_id'})
           .eqJoin('setup_id', r.table('setups')).zip()
           .merge(function (setup) {
               return {
                   properties: r.table('setupproperties')
                       .getAll(setup('setup_id'), {index: 'setup_id'})
                       .coerceTo('array')
               }
           }).coerceTo('array'),

       input_samples: r.table('process2sample').getAll(process('id'), {index: 'process_id'})
           .filter({'direction': 'in'})
           .eqJoin('sample_id', r.table('samples')).zip()
           .merge(function (sample) {
               return {
                   properties: r.table('propertyset2property')
                       .getAll(sample('property_set_id'), {index: 'property_set_id'})
                       .eqJoin('property_id', r.table('properties')).zip()
                       .orderBy('name')
                       .merge(function (property) {
                           return {
                               best_measure: r.table('best_measure_history')
                                   .getAll(property('best_measure_id'))
                                   .eqJoin('measurement_id', r.table('measurements'))
                                   .zip().coerceTo('array')
                           }
                       }).coerceTo('array'),
                   files: r.table('sample2datafile').getAll(sample('id'), {index: 'sample_id'})
                       .eqJoin('datafile_id', r.table('datafiles')).zip().pluck('id', 'name')
                       .coerceTo('array')
                   // processes: r.table('process2sample').getAll(sample('id'), {index: 'sample_id'})
                   //     .pluck('process_id', 'sample_id').distinct()
                   //     .eqJoin('process_id', r.table('processes')).zip().coerceTo('array')
               }
           }).coerceTo('array'),
       output_samples: r.table('process2sample').getAll(process('id'), {index: 'process_id'})
           .filter({'direction': 'out'})
           .eqJoin('sample_id', r.table('samples')).zip()
           .merge(function (sample) {
               return {
                   properties: r.table('propertyset2property')
                       .getAll(sample('property_set_id'), {index: 'property_set_id'})
                       .eqJoin('property_id', r.table('properties')).zip()
                       .orderBy('name')
                       .merge(function (property) {
                           return {
                               best_measure: r.table('best_measure_history')
                                   .getAll(property('best_measure_id'))
                                   .eqJoin('measurement_id', r.table('measurements'))
                                   .zip().coerceTo('array')
                           }
                       }).coerceTo('array'),
                   files: r.table('sample2datafile').getAll(sample('id'), {index: 'sample_id'})
                       .eqJoin('datafile_id', r.table('datafiles')).zip().pluck('id', 'name')
                       .coerceTo('array')
                   // processes: r.table('process2sample').getAll(sample('id'), {index: 'sample_id'})
                   //     .pluck('process_id', 'sample_id').distinct()
                   //     .eqJoin('process_id', r.table('processes')).zip().coerceTo('array')
               }
           }).coerceTo('array'),
       measurements: r.table('process2measurement').getAll(process('id'), {index: 'process_id'})
           .eqJoin('measurement_id', r.table('measurements')).zip()
           .merge(p2m => {
               return {
                   is_best_measure: r.db('materialscommons').table('best_measure_history')
                       .getAll(p2m('measurement_id'), {index: 'measurement_id'}).count()
               }
           })
           .coerceTo('array'),
       files_count: r.table('process2file').getAll(process('id'), {index: 'process_id'}).count(),
       files: [],
       filesLoaded: false,
       // files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
       //     .eqJoin('datafile_id', r.table('datafiles')).zip()
       //     .merge(f => {
       //         return {
       //             samples: r.table('sample2datafile').getAll(f('id'), {index: 'datafile_id'})
       //                 .eqJoin('sample_id', r.table('samples')).zip()
       //                 .distinct().coerceTo('array')
       //         };
       //     })
       //     .coerceTo('array'),
       input_files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
           .filter({direction: 'in'})
           .eqJoin('datafile_id', r.table('datafiles'))
           .zip().coerceTo('array'),
       output_files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
           .filter({direction: 'out'})
           .eqJoin('datafile_id', r.table('datafiles'))
           .zip().coerceTo('array')
   }
*/
